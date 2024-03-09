#include <fmt/core.h>

#include <fstream>
#include <future>
#include <iostream>
#include <map>
#include <sstream>
#include <vector>

int chunk_count = 0;

struct CityTemperature {
  float max_temperature;
  float min_temperature;
  float avg_temperature;
  float total_temperature;
  float temperature_count;
};

std::map<std::string, float> process_chunk(const std::string& chunk) {
  std::istringstream iss(chunk);
  std::string line;
  std::map<std::string, float> my_map;

  while (std::getline(iss, line, '\n')) {
    if (line.empty()) {
      continue;
    }

    // Assuming each line has the format "city;temperature"
    size_t delimiter_pos = line.find(';');
    if (delimiter_pos != std::string::npos) {
      std::string city = line.substr(0, delimiter_pos);
      // fmt::print("City: {}\n", city);
      float temperature;
      try {
        temperature = std::stof(line.substr(delimiter_pos + 1));
      } catch (std::invalid_argument& e) {
        fmt::print("Invalid temperature: {}\n", city);
        fmt::print("Invalid temperature: {}\n", line.substr(delimiter_pos + 1));
        continue;
      }
      my_map[city] += temperature;
    }
  }
  return my_map;
}

int main() {
  const std::string file_path = "../1bill.txt";

  std::ifstream file(file_path);
  if (!file.is_open()) {
    fmt::print("Failed to open file: {}\n", file_path);
    return 1;
  }
  const size_t buffer_size = 64 * 1024 * 1024;  // 64MB

  std::vector<char> buffer(buffer_size);

  std::vector<std::future<std::map<std::string, float>>> futures;

  while (file) {
    file.read(buffer.data(), buffer.size());
    std::streamsize bytesRead = file.gcount();
    std::string chunk(buffer.data(), bytesRead);

    futures.emplace_back(std::async(
        std::launch::async, [chunk]() { return process_chunk(chunk); }));
  }

  std::map<std::string, CityTemperature> final_map = {};

  for (auto& future : futures) {
    std::map<std::string, float> m = future.get();

    for (const auto& [city, temperature] : m) {
      if (final_map.find(city) != final_map.end()) {
        final_map[city].max_temperature =
            std::max(final_map[city].max_temperature, temperature);
        final_map[city].min_temperature =
            std::min(final_map[city].min_temperature, temperature);
        final_map[city].total_temperature += temperature;
        final_map[city].temperature_count++;
        final_map[city].avg_temperature =
            final_map[city].total_temperature / chunk_count;
      } else {
        final_map[city] = {temperature, temperature, temperature, temperature};
      }
    }
  }
  for (const auto& [city, temperature] : final_map) {
    // fmt::println("############################################");
    // fmt::print("City: {}\n", city);
    // fmt::print("Max temperature: {}\n", temperature.max_temperature);
    // fmt::print("Min temperature: {}\n", temperature.min_temperature);
    // fmt::print("Total temperature: {}\n", temperature.total_temperature);
    // fmt::print("Avg temperature: {}\n",
    //            temperature.total_temperature /
    //            temperature.temperature_count);
  }
  fmt::print("Done\n");
  return 0;
}
