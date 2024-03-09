
do:
	echo "building" && cmake -B build -S . -DCMAKE_TOOLCHAIN_FILE=/home/slim/repos/saleem/vcpkg/scripts/buildsystems/vcpkg.cmake -DCMAKE_EXPORT_COMPILE_COMMANDS=1 && cd build && make && cd -

buildbin:
	cd build && cmake .. && make && cd -

open:
	docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres
