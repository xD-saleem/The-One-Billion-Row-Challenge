cmake_minimum_required(VERSION 3.0)
project(main)

find_package(fmt REQUIRED)

# Set the path to parallel-hashmap headers
find_path(PARALLEL_HASHMAP_INCLUDE_DIRS "parallel_hashmap/btree.h")

# Check if parallel-hashmap headers are found
if(PARALLEL_HASHMAP_INCLUDE_DIRS)
    message(STATUS "parallel-hashmap headers found: ${PARALLEL_HASHMAP_INCLUDE_DIRS}")
else()
    message(FATAL_ERROR "parallel-hashmap headers not found. Please check the installation.")
endif()

# Include the parallel-hashmap headers
include_directories(${PARALLEL_HASHMAP_INCLUDE_DIRS})

add_executable(main main.cpp)
target_link_libraries(main PUBLIC fmt::fmt)

## replace with your vcpkg path
set(CMAKE_TOOLCHAIN_FILE "/home/slim/vcpkg/scripts/buildsystems/vcpkg.cmake"
    CACHE STRING "Vcpkg toolchain file")

set(CMAKE_BUILD_TYPE Debug)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -g")
