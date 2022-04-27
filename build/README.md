# About cadvisor

if `./cadvisor: /lib/x86_64-linux-gnu/libc.so.6: version GLIBC_2.32 not found (required by ./cadvisor)`

run following commands
````shell
wget http://ftp.gnu.org/gnu/glibc/glibc-2.23.tar.gz
tar -zxvf  glibc-2.32.tar.gz
mkdir glibc-2.32-build
cd glibc-2.32-build
../glibc-2.32/configure  --prefix=/usr --disable-profile --enable-add-ons --with-headers=/usr/include --with-binutils=/usr/bin
make -j8
make install
````

Packaging and Continuous Integration.

Put your cloud (AMI), container (Docker), OS (deb, rpm, pkg) package configurations and scripts in the `/build/package` directory.

Put your CI (travis, circle, drone) configurations and scripts in the `/build/ci` directory. Note that some of the CI tools (e.g., Travis CI) are very picky about the location of their config files. Try putting the config files in the `/build/ci` directory linking them to the location where the CI tools expect them when possible (don't worry if it's not and if keeping those files in the root directory makes your life easier :-)).

Examples:

* https://github.com/cockroachdb/cockroach/tree/master/build
