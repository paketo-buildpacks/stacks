package packages_test

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/image"
	"github.com/paketo-buildpacks/stacks/create-stack/pkg/packages"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBionic(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		buildImageTag = "paketobuildpacks/build:1.0.22-base-cnb"
		runImageTag   = "paketobuildpacks/run:1.0.22-base-cnb"
		bionic        packages.Bionic
	)

	it.Before(func() {
		client := image.Client{}

		_, err := client.Pull(buildImageTag, authn.DefaultKeychain)
		Expect(err).NotTo(HaveOccurred())

		_, err = client.Pull(runImageTag, authn.DefaultKeychain)
		Expect(err).NotTo(HaveOccurred())
	})

	it("can get the list of packages", func() {
		buildPackages, err := bionic.GetBuildPackagesList(buildImageTag)
		Expect(err).NotTo(HaveOccurred())
		Expect(buildPackages).To(Equal([]string{
			"adduser",
			"apt",
			"base-files",
			"base-passwd",
			"bash",
			"binutils",
			"binutils-common",
			"binutils-x86-64-linux-gnu",
			"bsdutils",
			"build-essential",
			"bzip2",
			"ca-certificates",
			"coreutils",
			"cpp",
			"cpp-7",
			"curl",
			"dash",
			"debconf",
			"debianutils",
			"diffutils",
			"dpkg",
			"dpkg-dev",
			"e2fsprogs",
			"fdisk",
			"findutils",
			"g++",
			"g++-7",
			"gcc",
			"gcc-7",
			"gcc-7-base",
			"gcc-8-base",
			"git",
			"git-man",
			"gpgv",
			"grep",
			"gzip",
			"hostname",
			"init-system-helpers",
			"jq",
			"libacl1",
			"libapt-pkg5.0",
			"libasan4",
			"libasn1-8-heimdal",
			"libatomic1",
			"libattr1",
			"libaudit-common",
			"libaudit1",
			"libbinutils",
			"libblkid1",
			"libbz2-1.0",
			"libc-bin",
			"libc-dev-bin",
			"libc6",
			"libc6-dev",
			"libcap-ng0",
			"libcc1-0",
			"libcilkrts5",
			"libcom-err2",
			"libcurl3-gnutls",
			"libcurl4",
			"libdb5.3",
			"libdebconfclient0",
			"libdpkg-perl",
			"liberror-perl",
			"libexpat1",
			"libext2fs2",
			"libfdisk1",
			"libffi6",
			"libgcc-7-dev",
			"libgcc1",
			"libgcrypt20",
			"libgdbm-compat4",
			"libgdbm5",
			"libgmp-dev",
			"libgmp10",
			"libgmpxx4ldbl",
			"libgnutls30",
			"libgomp1",
			"libgpg-error0",
			"libgssapi-krb5-2",
			"libgssapi3-heimdal",
			"libhcrypto4-heimdal",
			"libheimbase1-heimdal",
			"libheimntlm0-heimdal",
			"libhogweed4",
			"libhx509-5-heimdal",
			"libidn2-0",
			"libisl19",
			"libitm1",
			"libjq1",
			"libk5crypto3",
			"libkeyutils1",
			"libkrb5-26-heimdal",
			"libkrb5-3",
			"libkrb5support0",
			"libldap-2.4-2",
			"libldap-common",
			"liblsan0",
			"liblz4-1",
			"liblzma5",
			"libmount1",
			"libmpc3",
			"libmpfr6",
			"libmpx2",
			"libncurses5",
			"libncursesw5",
			"libnettle6",
			"libnghttp2-14",
			"libonig4",
			"libp11-kit0",
			"libpam-modules",
			"libpam-modules-bin",
			"libpam-runtime",
			"libpam0g",
			"libpcre3",
			"libperl5.26",
			"libprocps6",
			"libpsl5",
			"libquadmath0",
			"libroken18-heimdal",
			"librtmp1",
			"libsasl2-2",
			"libsasl2-modules-db",
			"libseccomp2",
			"libselinux1",
			"libsemanage-common",
			"libsemanage1",
			"libsepol1",
			"libsmartcols1",
			"libsqlite3-0",
			"libss2",
			"libssl1.1",
			"libstdc++-7-dev",
			"libstdc++6",
			"libsystemd0",
			"libtasn1-6",
			"libtinfo5",
			"libtsan0",
			"libubsan0",
			"libudev1",
			"libunistring2",
			"libuuid1",
			"libwind0-heimdal",
			"libyaml-0-2",
			"libzstd1",
			"linux-libc-dev",
			"locales",
			"login",
			"lsb-base",
			"make",
			"mawk",
			"mount",
			"ncurses-base",
			"ncurses-bin",
			"openssl",
			"passwd",
			"patch",
			"perl",
			"perl-base",
			"perl-modules-5.26",
			"procps",
			"sed",
			"sensible-utils",
			"sysvinit-utils",
			"tar",
			"tzdata",
			"ubuntu-keyring",
			"util-linux",
			"xz-utils",
			"zlib1g",
			"zlib1g-dev",
		}))

		runPackages, err := bionic.GetRunPackagesList(runImageTag)
		Expect(err).NotTo(HaveOccurred())
		Expect(runPackages).To(Equal([]string{
			"adduser",
			"apt",
			"base-files",
			"base-passwd",
			"bash",
			"bsdutils",
			"bzip2",
			"ca-certificates",
			"coreutils",
			"dash",
			"debconf",
			"debianutils",
			"diffutils",
			"dpkg",
			"e2fsprogs",
			"fdisk",
			"findutils",
			"gcc-8-base",
			"gpgv",
			"grep",
			"gzip",
			"hostname",
			"init-system-helpers",
			"libacl1",
			"libapt-pkg5.0",
			"libattr1",
			"libaudit-common",
			"libaudit1",
			"libblkid1",
			"libbz2-1.0",
			"libc-bin",
			"libc6",
			"libcap-ng0",
			"libcom-err2",
			"libdb5.3",
			"libdebconfclient0",
			"libext2fs2",
			"libfdisk1",
			"libffi6",
			"libgcc1",
			"libgcrypt20",
			"libgmp10",
			"libgnutls30",
			"libgpg-error0",
			"libhogweed4",
			"libidn2-0",
			"liblz4-1",
			"liblzma5",
			"libmount1",
			"libncurses5",
			"libncursesw5",
			"libnettle6",
			"libp11-kit0",
			"libpam-modules",
			"libpam-modules-bin",
			"libpam-runtime",
			"libpam0g",
			"libpcre3",
			"libprocps6",
			"libseccomp2",
			"libselinux1",
			"libsemanage-common",
			"libsemanage1",
			"libsepol1",
			"libsmartcols1",
			"libss2",
			"libssl1.1",
			"libstdc++6",
			"libsystemd0",
			"libtasn1-6",
			"libtinfo5",
			"libudev1",
			"libunistring2",
			"libuuid1",
			"libyaml-0-2",
			"libzstd1",
			"locales",
			"login",
			"lsb-base",
			"mawk",
			"mount",
			"ncurses-base",
			"ncurses-bin",
			"openssl",
			"passwd",
			"perl-base",
			"procps",
			"sed",
			"sensible-utils",
			"sysvinit-utils",
			"tar",
			"tzdata",
			"ubuntu-keyring",
			"util-linux",
			"zlib1g",
		}))
	})

	it("can get the package metadata", func() {
		buildPackageMetadata, err := bionic.GetBuildPackageMetadata(buildImageTag)
		Expect(err).NotTo(HaveOccurred())
		Expect(buildPackageMetadata).To(MatchJSON(`[
			{
				"arch": "all",
				"name": "adduser",
				"source": {
					"name": "adduser",
					"upstreamVersion": "3.116ubuntu1",
					"version": "3.116ubuntu1"
				},
				"summary": "add and remove users and groups",
				"version": "3.116ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "apt",
				"source": {
					"name": "apt",
					"upstreamVersion": "1.6.12ubuntu0.2",
					"version": "1.6.12ubuntu0.2"
				},
				"summary": "commandline package manager",
				"version": "1.6.12ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "base-files",
				"source": {
					"name": "base-files",
					"upstreamVersion": "10.1ubuntu2.10",
					"version": "10.1ubuntu2.10"
				},
				"summary": "Debian base system miscellaneous files",
				"version": "10.1ubuntu2.10"
			},
			{
				"arch": "amd64",
				"name": "base-passwd",
				"source": {
					"name": "base-passwd",
					"upstreamVersion": "3.5.44",
					"version": "3.5.44"
				},
				"summary": "Debian base system master password and group files",
				"version": "3.5.44"
			},
			{
				"arch": "amd64",
				"name": "bash",
				"source": {
					"name": "bash",
					"upstreamVersion": "4.4.18",
					"version": "4.4.18-2ubuntu1.2"
				},
				"summary": "GNU Bourne Again SHell",
				"version": "4.4.18-2ubuntu1.2"
			},
			{
				"arch": "amd64",
				"name": "binutils",
				"source": {
					"name": "binutils",
					"upstreamVersion": "2.30",
					"version": "2.30-21ubuntu1~18.04.4"
				},
				"summary": "GNU assembler, linker and binary utilities",
				"version": "2.30-21ubuntu1~18.04.4"
			},
			{
				"arch": "amd64",
				"name": "binutils-common:amd64",
				"source": {
					"name": "binutils",
					"upstreamVersion": "2.30",
					"version": "2.30-21ubuntu1~18.04.4"
				},
				"summary": "Common files for the GNU assembler, linker and binary utilities",
				"version": "2.30-21ubuntu1~18.04.4"
			},
			{
				"arch": "amd64",
				"name": "binutils-x86-64-linux-gnu",
				"source": {
					"name": "binutils",
					"upstreamVersion": "2.30",
					"version": "2.30-21ubuntu1~18.04.4"
				},
				"summary": "GNU binary utilities, for x86-64-linux-gnu target",
				"version": "2.30-21ubuntu1~18.04.4"
			},
			{
				"arch": "amd64",
				"name": "bsdutils",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "basic utilities from 4.4BSD-Lite",
				"version": "1:2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "build-essential",
				"source": {
					"name": "build-essential",
					"upstreamVersion": "12.4ubuntu1",
					"version": "12.4ubuntu1"
				},
				"summary": "Informational list of build-essential packages",
				"version": "12.4ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "bzip2",
				"source": {
					"name": "bzip2",
					"upstreamVersion": "1.0.6",
					"version": "1.0.6-8.1ubuntu0.2"
				},
				"summary": "high-quality block-sorting file compressor - utilities",
				"version": "1.0.6-8.1ubuntu0.2"
			},
			{
				"arch": "all",
				"name": "ca-certificates",
				"source": {
					"name": "ca-certificates",
					"upstreamVersion": "20210119~18.04.1",
					"version": "20210119~18.04.1"
				},
				"summary": "Common CA certificates",
				"version": "20210119~18.04.1"
			},
			{
				"arch": "amd64",
				"name": "coreutils",
				"source": {
					"name": "coreutils",
					"upstreamVersion": "8.28",
					"version": "8.28-1ubuntu1"
				},
				"summary": "GNU core utilities",
				"version": "8.28-1ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "cpp",
				"source": {
					"name": "gcc-defaults",
					"upstreamVersion": "1.176ubuntu2.3",
					"version": "1.176ubuntu2.3"
				},
				"summary": "GNU C preprocessor (cpp)",
				"version": "4:7.4.0-1ubuntu2.3"
			},
			{
				"arch": "amd64",
				"name": "cpp-7",
				"source": {
					"name": "gcc-7",
					"upstreamVersion": "7.5.0",
					"version": "7.5.0-3ubuntu1~18.04"
				},
				"summary": "GNU C preprocessor",
				"version": "7.5.0-3ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "curl",
				"source": {
					"name": "curl",
					"upstreamVersion": "7.58.0",
					"version": "7.58.0-2ubuntu3.12"
				},
				"summary": "command line tool for transferring data with URL syntax",
				"version": "7.58.0-2ubuntu3.12"
			},
			{
				"arch": "amd64",
				"name": "dash",
				"source": {
					"name": "dash",
					"upstreamVersion": "0.5.8",
					"version": "0.5.8-2.10"
				},
				"summary": "POSIX-compliant shell",
				"version": "0.5.8-2.10"
			},
			{
				"arch": "all",
				"name": "debconf",
				"source": {
					"name": "debconf",
					"upstreamVersion": "1.5.66ubuntu1",
					"version": "1.5.66ubuntu1"
				},
				"summary": "Debian configuration management system",
				"version": "1.5.66ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "debianutils",
				"source": {
					"name": "debianutils",
					"upstreamVersion": "4.8.4",
					"version": "4.8.4"
				},
				"summary": "Miscellaneous utilities specific to Debian",
				"version": "4.8.4"
			},
			{
				"arch": "amd64",
				"name": "diffutils",
				"source": {
					"name": "diffutils",
					"upstreamVersion": "3.6",
					"version": "1:3.6-1"
				},
				"summary": "File comparison utilities",
				"version": "1:3.6-1"
			},
			{
				"arch": "amd64",
				"name": "dpkg",
				"source": {
					"name": "dpkg",
					"upstreamVersion": "1.19.0.5ubuntu2.3",
					"version": "1.19.0.5ubuntu2.3"
				},
				"summary": "Debian package management system",
				"version": "1.19.0.5ubuntu2.3"
			},
			{
				"arch": "all",
				"name": "dpkg-dev",
				"source": {
					"name": "dpkg",
					"upstreamVersion": "1.19.0.5ubuntu2.3",
					"version": "1.19.0.5ubuntu2.3"
				},
				"summary": "Debian package development tools",
				"version": "1.19.0.5ubuntu2.3"
			},
			{
				"arch": "amd64",
				"name": "e2fsprogs",
				"source": {
					"name": "e2fsprogs",
					"upstreamVersion": "1.44.1",
					"version": "1.44.1-1ubuntu1.3"
				},
				"summary": "ext2/ext3/ext4 file system utilities",
				"version": "1.44.1-1ubuntu1.3"
			},
			{
				"arch": "amd64",
				"name": "fdisk",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "collection of partitioning utilities",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "findutils",
				"source": {
					"name": "findutils",
					"upstreamVersion": "4.6.0+git+20170828",
					"version": "4.6.0+git+20170828-2"
				},
				"summary": "utilities for finding files--find, xargs",
				"version": "4.6.0+git+20170828-2"
			},
			{
				"arch": "amd64",
				"name": "g++",
				"source": {
					"name": "gcc-defaults",
					"upstreamVersion": "1.176ubuntu2.3",
					"version": "1.176ubuntu2.3"
				},
				"summary": "GNU C++ compiler",
				"version": "4:7.4.0-1ubuntu2.3"
			},
			{
				"arch": "amd64",
				"name": "g++-7",
				"source": {
					"name": "gcc-7",
					"upstreamVersion": "7.5.0",
					"version": "7.5.0-3ubuntu1~18.04"
				},
				"summary": "GNU C++ compiler",
				"version": "7.5.0-3ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "gcc",
				"source": {
					"name": "gcc-defaults",
					"upstreamVersion": "1.176ubuntu2.3",
					"version": "1.176ubuntu2.3"
				},
				"summary": "GNU C compiler",
				"version": "4:7.4.0-1ubuntu2.3"
			},
			{
				"arch": "amd64",
				"name": "gcc-7",
				"source": {
					"name": "gcc-7",
					"upstreamVersion": "7.5.0",
					"version": "7.5.0-3ubuntu1~18.04"
				},
				"summary": "GNU C compiler",
				"version": "7.5.0-3ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "gcc-7-base:amd64",
				"source": {
					"name": "gcc-7",
					"upstreamVersion": "7.5.0",
					"version": "7.5.0-3ubuntu1~18.04"
				},
				"summary": "GCC, the GNU Compiler Collection (base package)",
				"version": "7.5.0-3ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "gcc-8-base:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GCC, the GNU Compiler Collection (base package)",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "git",
				"source": {
					"name": "git",
					"upstreamVersion": "2.17.1",
					"version": "1:2.17.1-1ubuntu0.7"
				},
				"summary": "fast, scalable, distributed revision control system",
				"version": "1:2.17.1-1ubuntu0.7"
			},
			{
				"arch": "all",
				"name": "git-man",
				"source": {
					"name": "git",
					"upstreamVersion": "2.17.1",
					"version": "1:2.17.1-1ubuntu0.7"
				},
				"summary": "fast, scalable, distributed revision control system (manual pages)",
				"version": "1:2.17.1-1ubuntu0.7"
			},
			{
				"arch": "amd64",
				"name": "gpgv",
				"source": {
					"name": "gnupg2",
					"upstreamVersion": "2.2.4",
					"version": "2.2.4-1ubuntu1.4"
				},
				"summary": "GNU privacy guard - signature verification tool",
				"version": "2.2.4-1ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "grep",
				"source": {
					"name": "grep",
					"upstreamVersion": "3.1",
					"version": "3.1-2build1"
				},
				"summary": "GNU grep, egrep and fgrep",
				"version": "3.1-2build1"
			},
			{
				"arch": "amd64",
				"name": "gzip",
				"source": {
					"name": "gzip",
					"upstreamVersion": "1.6",
					"version": "1.6-5ubuntu1"
				},
				"summary": "GNU compression utilities",
				"version": "1.6-5ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "hostname",
				"source": {
					"name": "hostname",
					"upstreamVersion": "3.20",
					"version": "3.20"
				},
				"summary": "utility to set/show the host name or domain name",
				"version": "3.20"
			},
			{
				"arch": "all",
				"name": "init-system-helpers",
				"source": {
					"name": "init-system-helpers",
					"upstreamVersion": "1.51",
					"version": "1.51"
				},
				"summary": "helper tools for all init systems",
				"version": "1.51"
			},
			{
				"arch": "amd64",
				"name": "jq",
				"source": {
					"name": "jq",
					"upstreamVersion": "1.5+dfsg",
					"version": "1.5+dfsg-2"
				},
				"summary": "lightweight and flexible command-line JSON processor",
				"version": "1.5+dfsg-2"
			},
			{
				"arch": "amd64",
				"name": "libacl1:amd64",
				"source": {
					"name": "acl",
					"upstreamVersion": "2.2.52",
					"version": "2.2.52-3build1"
				},
				"summary": "Access control list shared library",
				"version": "2.2.52-3build1"
			},
			{
				"arch": "amd64",
				"name": "libapt-pkg5.0:amd64",
				"source": {
					"name": "apt",
					"upstreamVersion": "1.6.12ubuntu0.2",
					"version": "1.6.12ubuntu0.2"
				},
				"summary": "package management runtime library",
				"version": "1.6.12ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libasan4:amd64",
				"source": {
					"name": "gcc-7",
					"upstreamVersion": "7.5.0",
					"version": "7.5.0-3ubuntu1~18.04"
				},
				"summary": "AddressSanitizer -- a fast memory error detector",
				"version": "7.5.0-3ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libasn1-8-heimdal:amd64",
				"source": {
					"name": "heimdal",
					"upstreamVersion": "7.5.0+dfsg",
					"version": "7.5.0+dfsg-1"
				},
				"summary": "Heimdal Kerberos - ASN.1 library",
				"version": "7.5.0+dfsg-1"
			},
			{
				"arch": "amd64",
				"name": "libatomic1:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "support library providing __atomic built-in functions",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libattr1:amd64",
				"source": {
					"name": "attr",
					"upstreamVersion": "2.4.47",
					"version": "1:2.4.47-2build1"
				},
				"summary": "Extended attribute shared library",
				"version": "1:2.4.47-2build1"
			},
			{
				"arch": "all",
				"name": "libaudit-common",
				"source": {
					"name": "audit",
					"upstreamVersion": "2.8.2",
					"version": "1:2.8.2-1ubuntu1.1"
				},
				"summary": "Dynamic library for security auditing - common files",
				"version": "1:2.8.2-1ubuntu1.1"
			},
			{
				"arch": "amd64",
				"name": "libaudit1:amd64",
				"source": {
					"name": "audit",
					"upstreamVersion": "2.8.2",
					"version": "1:2.8.2-1ubuntu1.1"
				},
				"summary": "Dynamic library for security auditing",
				"version": "1:2.8.2-1ubuntu1.1"
			},
			{
				"arch": "amd64",
				"name": "libbinutils:amd64",
				"source": {
					"name": "binutils",
					"upstreamVersion": "2.30",
					"version": "2.30-21ubuntu1~18.04.4"
				},
				"summary": "GNU binary utilities (private shared library)",
				"version": "2.30-21ubuntu1~18.04.4"
			},
			{
				"arch": "amd64",
				"name": "libblkid1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "block device ID library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libbz2-1.0:amd64",
				"source": {
					"name": "bzip2",
					"upstreamVersion": "1.0.6",
					"version": "1.0.6-8.1ubuntu0.2"
				},
				"summary": "high-quality block-sorting file compressor library - runtime",
				"version": "1.0.6-8.1ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libc-bin",
				"source": {
					"name": "glibc",
					"upstreamVersion": "2.27",
					"version": "2.27-3ubuntu1.4"
				},
				"summary": "GNU C Library: Binaries",
				"version": "2.27-3ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "libc-dev-bin",
				"source": {
					"name": "glibc",
					"upstreamVersion": "2.27",
					"version": "2.27-3ubuntu1.4"
				},
				"summary": "GNU C Library: Development binaries",
				"version": "2.27-3ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "libc6:amd64",
				"source": {
					"name": "glibc",
					"upstreamVersion": "2.27",
					"version": "2.27-3ubuntu1.4"
				},
				"summary": "GNU C Library: Shared libraries",
				"version": "2.27-3ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "libc6-dev:amd64",
				"source": {
					"name": "glibc",
					"upstreamVersion": "2.27",
					"version": "2.27-3ubuntu1.4"
				},
				"summary": "GNU C Library: Development Libraries and Header Files",
				"version": "2.27-3ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "libcap-ng0:amd64",
				"source": {
					"name": "libcap-ng",
					"upstreamVersion": "0.7.7",
					"version": "0.7.7-3.1"
				},
				"summary": "An alternate POSIX capabilities library",
				"version": "0.7.7-3.1"
			},
			{
				"arch": "amd64",
				"name": "libcc1-0:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GCC cc1 plugin for GDB",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libcilkrts5:amd64",
				"source": {
					"name": "gcc-7",
					"upstreamVersion": "7.5.0",
					"version": "7.5.0-3ubuntu1~18.04"
				},
				"summary": "Intel Cilk Plus language extensions (runtime)",
				"version": "7.5.0-3ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libcom-err2:amd64",
				"source": {
					"name": "e2fsprogs",
					"upstreamVersion": "1.44.1",
					"version": "1.44.1-1ubuntu1.3"
				},
				"summary": "common error description library",
				"version": "1.44.1-1ubuntu1.3"
			},
			{
				"arch": "amd64",
				"name": "libcurl3-gnutls:amd64",
				"source": {
					"name": "curl",
					"upstreamVersion": "7.58.0",
					"version": "7.58.0-2ubuntu3.12"
				},
				"summary": "easy-to-use client-side URL transfer library (GnuTLS flavour)",
				"version": "7.58.0-2ubuntu3.12"
			},
			{
				"arch": "amd64",
				"name": "libcurl4:amd64",
				"source": {
					"name": "curl",
					"upstreamVersion": "7.58.0",
					"version": "7.58.0-2ubuntu3.12"
				},
				"summary": "easy-to-use client-side URL transfer library (OpenSSL flavour)",
				"version": "7.58.0-2ubuntu3.12"
			},
			{
				"arch": "amd64",
				"name": "libdb5.3:amd64",
				"source": {
					"name": "db5.3",
					"upstreamVersion": "5.3.28",
					"version": "5.3.28-13.1ubuntu1.1"
				},
				"summary": "Berkeley v5.3 Database Libraries [runtime]",
				"version": "5.3.28-13.1ubuntu1.1"
			},
			{
				"arch": "amd64",
				"name": "libdebconfclient0:amd64",
				"source": {
					"name": "cdebconf",
					"upstreamVersion": "0.213ubuntu1",
					"version": "0.213ubuntu1"
				},
				"summary": "Debian Configuration Management System (C-implementation library)",
				"version": "0.213ubuntu1"
			},
			{
				"arch": "all",
				"name": "libdpkg-perl",
				"source": {
					"name": "dpkg",
					"upstreamVersion": "1.19.0.5ubuntu2.3",
					"version": "1.19.0.5ubuntu2.3"
				},
				"summary": "Dpkg perl modules",
				"version": "1.19.0.5ubuntu2.3"
			},
			{
				"arch": "all",
				"name": "liberror-perl",
				"source": {
					"name": "liberror-perl",
					"upstreamVersion": "0.17025",
					"version": "0.17025-1"
				},
				"summary": "Perl module for error/exception handling in an OO-ish way",
				"version": "0.17025-1"
			},
			{
				"arch": "amd64",
				"name": "libexpat1:amd64",
				"source": {
					"name": "expat",
					"upstreamVersion": "2.2.5",
					"version": "2.2.5-3ubuntu0.2"
				},
				"summary": "XML parsing C library - runtime library",
				"version": "2.2.5-3ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libext2fs2:amd64",
				"source": {
					"name": "e2fsprogs",
					"upstreamVersion": "1.44.1",
					"version": "1.44.1-1ubuntu1.3"
				},
				"summary": "ext2/ext3/ext4 file system libraries",
				"version": "1.44.1-1ubuntu1.3"
			},
			{
				"arch": "amd64",
				"name": "libfdisk1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "fdisk partitioning library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libffi6:amd64",
				"source": {
					"name": "libffi",
					"upstreamVersion": "3.2.1",
					"version": "3.2.1-8"
				},
				"summary": "Foreign Function Interface library runtime",
				"version": "3.2.1-8"
			},
			{
				"arch": "amd64",
				"name": "libgcc-7-dev:amd64",
				"source": {
					"name": "gcc-7",
					"upstreamVersion": "7.5.0",
					"version": "7.5.0-3ubuntu1~18.04"
				},
				"summary": "GCC support library (development files)",
				"version": "7.5.0-3ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libgcc1:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GCC support library",
				"version": "1:8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libgcrypt20:amd64",
				"source": {
					"name": "libgcrypt20",
					"upstreamVersion": "1.8.1",
					"version": "1.8.1-4ubuntu1.2"
				},
				"summary": "LGPL Crypto library - runtime library",
				"version": "1.8.1-4ubuntu1.2"
			},
			{
				"arch": "amd64",
				"name": "libgdbm-compat4:amd64",
				"source": {
					"name": "gdbm",
					"upstreamVersion": "1.14.1",
					"version": "1.14.1-6"
				},
				"summary": "GNU dbm database routines (legacy support runtime version)",
				"version": "1.14.1-6"
			},
			{
				"arch": "amd64",
				"name": "libgdbm5:amd64",
				"source": {
					"name": "gdbm",
					"upstreamVersion": "1.14.1",
					"version": "1.14.1-6"
				},
				"summary": "GNU dbm database routines (runtime version)",
				"version": "1.14.1-6"
			},
			{
				"arch": "amd64",
				"name": "libgmp-dev:amd64",
				"source": {
					"name": "gmp",
					"upstreamVersion": "6.1.2+dfsg",
					"version": "2:6.1.2+dfsg-2"
				},
				"summary": "Multiprecision arithmetic library developers tools",
				"version": "2:6.1.2+dfsg-2"
			},
			{
				"arch": "amd64",
				"name": "libgmp10:amd64",
				"source": {
					"name": "gmp",
					"upstreamVersion": "6.1.2+dfsg",
					"version": "2:6.1.2+dfsg-2"
				},
				"summary": "Multiprecision arithmetic library",
				"version": "2:6.1.2+dfsg-2"
			},
			{
				"arch": "amd64",
				"name": "libgmpxx4ldbl:amd64",
				"source": {
					"name": "gmp",
					"upstreamVersion": "6.1.2+dfsg",
					"version": "2:6.1.2+dfsg-2"
				},
				"summary": "Multiprecision arithmetic library (C++ bindings)",
				"version": "2:6.1.2+dfsg-2"
			},
			{
				"arch": "amd64",
				"name": "libgnutls30:amd64",
				"source": {
					"name": "gnutls28",
					"upstreamVersion": "3.5.18",
					"version": "3.5.18-1ubuntu1.4"
				},
				"summary": "GNU TLS library - main runtime library",
				"version": "3.5.18-1ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "libgomp1:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GCC OpenMP (GOMP) support library",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libgpg-error0:amd64",
				"source": {
					"name": "libgpg-error",
					"upstreamVersion": "1.27",
					"version": "1.27-6"
				},
				"summary": "library for common error values and messages in GnuPG components",
				"version": "1.27-6"
			},
			{
				"arch": "amd64",
				"name": "libgssapi-krb5-2:amd64",
				"source": {
					"name": "krb5",
					"upstreamVersion": "1.16",
					"version": "1.16-2ubuntu0.2"
				},
				"summary": "MIT Kerberos runtime libraries - krb5 GSS-API Mechanism",
				"version": "1.16-2ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libgssapi3-heimdal:amd64",
				"source": {
					"name": "heimdal",
					"upstreamVersion": "7.5.0+dfsg",
					"version": "7.5.0+dfsg-1"
				},
				"summary": "Heimdal Kerberos - GSSAPI support library",
				"version": "7.5.0+dfsg-1"
			},
			{
				"arch": "amd64",
				"name": "libhcrypto4-heimdal:amd64",
				"source": {
					"name": "heimdal",
					"upstreamVersion": "7.5.0+dfsg",
					"version": "7.5.0+dfsg-1"
				},
				"summary": "Heimdal Kerberos - crypto library",
				"version": "7.5.0+dfsg-1"
			},
			{
				"arch": "amd64",
				"name": "libheimbase1-heimdal:amd64",
				"source": {
					"name": "heimdal",
					"upstreamVersion": "7.5.0+dfsg",
					"version": "7.5.0+dfsg-1"
				},
				"summary": "Heimdal Kerberos - Base library",
				"version": "7.5.0+dfsg-1"
			},
			{
				"arch": "amd64",
				"name": "libheimntlm0-heimdal:amd64",
				"source": {
					"name": "heimdal",
					"upstreamVersion": "7.5.0+dfsg",
					"version": "7.5.0+dfsg-1"
				},
				"summary": "Heimdal Kerberos - NTLM support library",
				"version": "7.5.0+dfsg-1"
			},
			{
				"arch": "amd64",
				"name": "libhogweed4:amd64",
				"source": {
					"name": "nettle",
					"upstreamVersion": "3.4",
					"version": "3.4-1"
				},
				"summary": "low level cryptographic library (public-key cryptos)",
				"version": "3.4-1"
			},
			{
				"arch": "amd64",
				"name": "libhx509-5-heimdal:amd64",
				"source": {
					"name": "heimdal",
					"upstreamVersion": "7.5.0+dfsg",
					"version": "7.5.0+dfsg-1"
				},
				"summary": "Heimdal Kerberos - X509 support library",
				"version": "7.5.0+dfsg-1"
			},
			{
				"arch": "amd64",
				"name": "libidn2-0:amd64",
				"source": {
					"name": "libidn2",
					"upstreamVersion": "2.0.4",
					"version": "2.0.4-1.1ubuntu0.2"
				},
				"summary": "Internationalized domain names (IDNA2008/TR46) library",
				"version": "2.0.4-1.1ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libisl19:amd64",
				"source": {
					"name": "isl",
					"upstreamVersion": "0.19",
					"version": "0.19-1"
				},
				"summary": "manipulating sets and relations of integer points bounded by linear constraints",
				"version": "0.19-1"
			},
			{
				"arch": "amd64",
				"name": "libitm1:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GNU Transactional Memory Library",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libjq1:amd64",
				"source": {
					"name": "jq",
					"upstreamVersion": "1.5+dfsg",
					"version": "1.5+dfsg-2"
				},
				"summary": "lightweight and flexible command-line JSON processor - shared library",
				"version": "1.5+dfsg-2"
			},
			{
				"arch": "amd64",
				"name": "libk5crypto3:amd64",
				"source": {
					"name": "krb5",
					"upstreamVersion": "1.16",
					"version": "1.16-2ubuntu0.2"
				},
				"summary": "MIT Kerberos runtime libraries - Crypto Library",
				"version": "1.16-2ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libkeyutils1:amd64",
				"source": {
					"name": "keyutils",
					"upstreamVersion": "1.5.9",
					"version": "1.5.9-9.2ubuntu2"
				},
				"summary": "Linux Key Management Utilities (library)",
				"version": "1.5.9-9.2ubuntu2"
			},
			{
				"arch": "amd64",
				"name": "libkrb5-26-heimdal:amd64",
				"source": {
					"name": "heimdal",
					"upstreamVersion": "7.5.0+dfsg",
					"version": "7.5.0+dfsg-1"
				},
				"summary": "Heimdal Kerberos - libraries",
				"version": "7.5.0+dfsg-1"
			},
			{
				"arch": "amd64",
				"name": "libkrb5-3:amd64",
				"source": {
					"name": "krb5",
					"upstreamVersion": "1.16",
					"version": "1.16-2ubuntu0.2"
				},
				"summary": "MIT Kerberos runtime libraries",
				"version": "1.16-2ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libkrb5support0:amd64",
				"source": {
					"name": "krb5",
					"upstreamVersion": "1.16",
					"version": "1.16-2ubuntu0.2"
				},
				"summary": "MIT Kerberos runtime libraries - Support library",
				"version": "1.16-2ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libldap-2.4-2:amd64",
				"source": {
					"name": "openldap",
					"upstreamVersion": "2.4.45+dfsg",
					"version": "2.4.45+dfsg-1ubuntu1.10"
				},
				"summary": "OpenLDAP libraries",
				"version": "2.4.45+dfsg-1ubuntu1.10"
			},
			{
				"arch": "all",
				"name": "libldap-common",
				"source": {
					"name": "openldap",
					"upstreamVersion": "2.4.45+dfsg",
					"version": "2.4.45+dfsg-1ubuntu1.10"
				},
				"summary": "OpenLDAP common files for libraries",
				"version": "2.4.45+dfsg-1ubuntu1.10"
			},
			{
				"arch": "amd64",
				"name": "liblsan0:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "LeakSanitizer -- a memory leak detector (runtime)",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "liblz4-1:amd64",
				"source": {
					"name": "lz4",
					"upstreamVersion": "0.0~r131",
					"version": "0.0~r131-2ubuntu3"
				},
				"summary": "Fast LZ compression algorithm library - runtime",
				"version": "0.0~r131-2ubuntu3"
			},
			{
				"arch": "amd64",
				"name": "liblzma5:amd64",
				"source": {
					"name": "xz-utils",
					"upstreamVersion": "5.2.2",
					"version": "5.2.2-1.3"
				},
				"summary": "XZ-format compression library",
				"version": "5.2.2-1.3"
			},
			{
				"arch": "amd64",
				"name": "libmount1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "device mounting library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libmpc3:amd64",
				"source": {
					"name": "mpclib3",
					"upstreamVersion": "1.1.0",
					"version": "1.1.0-1"
				},
				"summary": "multiple precision complex floating-point library",
				"version": "1.1.0-1"
			},
			{
				"arch": "amd64",
				"name": "libmpfr6:amd64",
				"source": {
					"name": "mpfr4",
					"upstreamVersion": "4.0.1",
					"version": "4.0.1-1"
				},
				"summary": "multiple precision floating-point computation",
				"version": "4.0.1-1"
			},
			{
				"arch": "amd64",
				"name": "libmpx2:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "Intel memory protection extensions (runtime)",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libncurses5:amd64",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "shared libraries for terminal handling",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "libncursesw5:amd64",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "shared libraries for terminal handling (wide character support)",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "libnettle6:amd64",
				"source": {
					"name": "nettle",
					"upstreamVersion": "3.4",
					"version": "3.4-1"
				},
				"summary": "low level cryptographic library (symmetric and one-way cryptos)",
				"version": "3.4-1"
			},
			{
				"arch": "amd64",
				"name": "libnghttp2-14:amd64",
				"source": {
					"name": "nghttp2",
					"upstreamVersion": "1.30.0",
					"version": "1.30.0-1ubuntu1"
				},
				"summary": "library implementing HTTP/2 protocol (shared library)",
				"version": "1.30.0-1ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "libonig4:amd64",
				"source": {
					"name": "libonig",
					"upstreamVersion": "6.7.0",
					"version": "6.7.0-1"
				},
				"summary": "regular expressions library",
				"version": "6.7.0-1"
			},
			{
				"arch": "amd64",
				"name": "libp11-kit0:amd64",
				"source": {
					"name": "p11-kit",
					"upstreamVersion": "0.23.9",
					"version": "0.23.9-2ubuntu0.1"
				},
				"summary": "library for loading and coordinating access to PKCS#11 modules - runtime",
				"version": "0.23.9-2ubuntu0.1"
			},
			{
				"arch": "amd64",
				"name": "libpam-modules:amd64",
				"source": {
					"name": "pam",
					"upstreamVersion": "1.1.8",
					"version": "1.1.8-3.6ubuntu2.18.04.2"
				},
				"summary": "Pluggable Authentication Modules for PAM",
				"version": "1.1.8-3.6ubuntu2.18.04.2"
			},
			{
				"arch": "amd64",
				"name": "libpam-modules-bin",
				"source": {
					"name": "pam",
					"upstreamVersion": "1.1.8",
					"version": "1.1.8-3.6ubuntu2.18.04.2"
				},
				"summary": "Pluggable Authentication Modules for PAM - helper binaries",
				"version": "1.1.8-3.6ubuntu2.18.04.2"
			},
			{
				"arch": "all",
				"name": "libpam-runtime",
				"source": {
					"name": "pam",
					"upstreamVersion": "1.1.8",
					"version": "1.1.8-3.6ubuntu2.18.04.2"
				},
				"summary": "Runtime support for the PAM library",
				"version": "1.1.8-3.6ubuntu2.18.04.2"
			},
			{
				"arch": "amd64",
				"name": "libpam0g:amd64",
				"source": {
					"name": "pam",
					"upstreamVersion": "1.1.8",
					"version": "1.1.8-3.6ubuntu2.18.04.2"
				},
				"summary": "Pluggable Authentication Modules library",
				"version": "1.1.8-3.6ubuntu2.18.04.2"
			},
			{
				"arch": "amd64",
				"name": "libpcre3:amd64",
				"source": {
					"name": "pcre3",
					"upstreamVersion": "8.39",
					"version": "2:8.39-9"
				},
				"summary": "Old Perl 5 Compatible Regular Expression Library - runtime files",
				"version": "2:8.39-9"
			},
			{
				"arch": "amd64",
				"name": "libperl5.26:amd64",
				"source": {
					"name": "perl",
					"upstreamVersion": "5.26.1",
					"version": "5.26.1-6ubuntu0.5"
				},
				"summary": "shared Perl library",
				"version": "5.26.1-6ubuntu0.5"
			},
			{
				"arch": "amd64",
				"name": "libprocps6:amd64",
				"source": {
					"name": "procps",
					"upstreamVersion": "3.3.12",
					"version": "2:3.3.12-3ubuntu1.2"
				},
				"summary": "library for accessing process information from /proc",
				"version": "2:3.3.12-3ubuntu1.2"
			},
			{
				"arch": "amd64",
				"name": "libpsl5:amd64",
				"source": {
					"name": "libpsl",
					"upstreamVersion": "0.19.1",
					"version": "0.19.1-5build1"
				},
				"summary": "Library for Public Suffix List (shared libraries)",
				"version": "0.19.1-5build1"
			},
			{
				"arch": "amd64",
				"name": "libquadmath0:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GCC Quad-Precision Math Library",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libroken18-heimdal:amd64",
				"source": {
					"name": "heimdal",
					"upstreamVersion": "7.5.0+dfsg",
					"version": "7.5.0+dfsg-1"
				},
				"summary": "Heimdal Kerberos - roken support library",
				"version": "7.5.0+dfsg-1"
			},
			{
				"arch": "amd64",
				"name": "librtmp1:amd64",
				"source": {
					"name": "rtmpdump",
					"upstreamVersion": "2.4+20151223.gitfa8646d.1",
					"version": "2.4+20151223.gitfa8646d.1-1"
				},
				"summary": "toolkit for RTMP streams (shared library)",
				"version": "2.4+20151223.gitfa8646d.1-1"
			},
			{
				"arch": "amd64",
				"name": "libsasl2-2:amd64",
				"source": {
					"name": "cyrus-sasl2",
					"upstreamVersion": "2.1.27~101-g0780600+dfsg",
					"version": "2.1.27~101-g0780600+dfsg-3ubuntu2.3"
				},
				"summary": "Cyrus SASL - authentication abstraction library",
				"version": "2.1.27~101-g0780600+dfsg-3ubuntu2.3"
			},
			{
				"arch": "amd64",
				"name": "libsasl2-modules-db:amd64",
				"source": {
					"name": "cyrus-sasl2",
					"upstreamVersion": "2.1.27~101-g0780600+dfsg",
					"version": "2.1.27~101-g0780600+dfsg-3ubuntu2.3"
				},
				"summary": "Cyrus SASL - pluggable authentication modules (DB)",
				"version": "2.1.27~101-g0780600+dfsg-3ubuntu2.3"
			},
			{
				"arch": "amd64",
				"name": "libseccomp2:amd64",
				"source": {
					"name": "libseccomp",
					"upstreamVersion": "2.4.3",
					"version": "2.4.3-1ubuntu3.18.04.3"
				},
				"summary": "high level interface to Linux seccomp filter",
				"version": "2.4.3-1ubuntu3.18.04.3"
			},
			{
				"arch": "amd64",
				"name": "libselinux1:amd64",
				"source": {
					"name": "libselinux",
					"upstreamVersion": "2.7",
					"version": "2.7-2build2"
				},
				"summary": "SELinux runtime shared libraries",
				"version": "2.7-2build2"
			},
			{
				"arch": "all",
				"name": "libsemanage-common",
				"source": {
					"name": "libsemanage",
					"upstreamVersion": "2.7",
					"version": "2.7-2build2"
				},
				"summary": "Common files for SELinux policy management libraries",
				"version": "2.7-2build2"
			},
			{
				"arch": "amd64",
				"name": "libsemanage1:amd64",
				"source": {
					"name": "libsemanage",
					"upstreamVersion": "2.7",
					"version": "2.7-2build2"
				},
				"summary": "SELinux policy management library",
				"version": "2.7-2build2"
			},
			{
				"arch": "amd64",
				"name": "libsepol1:amd64",
				"source": {
					"name": "libsepol",
					"upstreamVersion": "2.7",
					"version": "2.7-1"
				},
				"summary": "SELinux library for manipulating binary security policies",
				"version": "2.7-1"
			},
			{
				"arch": "amd64",
				"name": "libsmartcols1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "smart column output alignment library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libsqlite3-0:amd64",
				"source": {
					"name": "sqlite3",
					"upstreamVersion": "3.22.0",
					"version": "3.22.0-1ubuntu0.4"
				},
				"summary": "SQLite 3 shared library",
				"version": "3.22.0-1ubuntu0.4"
			},
			{
				"arch": "amd64",
				"name": "libss2:amd64",
				"source": {
					"name": "e2fsprogs",
					"upstreamVersion": "1.44.1",
					"version": "1.44.1-1ubuntu1.3"
				},
				"summary": "command-line interface parsing library",
				"version": "1.44.1-1ubuntu1.3"
			},
			{
				"arch": "amd64",
				"name": "libssl1.1:amd64",
				"source": {
					"name": "openssl",
					"upstreamVersion": "1.1.1",
					"version": "1.1.1-1ubuntu2.1~18.04.8"
				},
				"summary": "Secure Sockets Layer toolkit - shared libraries",
				"version": "1.1.1-1ubuntu2.1~18.04.8"
			},
			{
				"arch": "amd64",
				"name": "libstdc++-7-dev:amd64",
				"source": {
					"name": "gcc-7",
					"upstreamVersion": "7.5.0",
					"version": "7.5.0-3ubuntu1~18.04"
				},
				"summary": "GNU Standard C++ Library v3 (development files)",
				"version": "7.5.0-3ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libstdc++6:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GNU Standard C++ Library v3",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libsystemd0:amd64",
				"source": {
					"name": "systemd",
					"upstreamVersion": "237",
					"version": "237-3ubuntu10.44"
				},
				"summary": "systemd utility library",
				"version": "237-3ubuntu10.44"
			},
			{
				"arch": "amd64",
				"name": "libtasn1-6:amd64",
				"source": {
					"name": "libtasn1-6",
					"upstreamVersion": "4.13",
					"version": "4.13-2"
				},
				"summary": "Manage ASN.1 structures (runtime)",
				"version": "4.13-2"
			},
			{
				"arch": "amd64",
				"name": "libtinfo5:amd64",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "shared low-level terminfo library for terminal handling",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "libtsan0:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "ThreadSanitizer -- a Valgrind-based detector of data races (runtime)",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libubsan0:amd64",
				"source": {
					"name": "gcc-7",
					"upstreamVersion": "7.5.0",
					"version": "7.5.0-3ubuntu1~18.04"
				},
				"summary": "UBSan -- undefined behaviour sanitizer (runtime)",
				"version": "7.5.0-3ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libudev1:amd64",
				"source": {
					"name": "systemd",
					"upstreamVersion": "237",
					"version": "237-3ubuntu10.44"
				},
				"summary": "libudev shared library",
				"version": "237-3ubuntu10.44"
			},
			{
				"arch": "amd64",
				"name": "libunistring2:amd64",
				"source": {
					"name": "libunistring",
					"upstreamVersion": "0.9.9",
					"version": "0.9.9-0ubuntu2"
				},
				"summary": "Unicode string library for C",
				"version": "0.9.9-0ubuntu2"
			},
			{
				"arch": "amd64",
				"name": "libuuid1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "Universally Unique ID library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libwind0-heimdal:amd64",
				"source": {
					"name": "heimdal",
					"upstreamVersion": "7.5.0+dfsg",
					"version": "7.5.0+dfsg-1"
				},
				"summary": "Heimdal Kerberos - stringprep implementation",
				"version": "7.5.0+dfsg-1"
			},
			{
				"arch": "amd64",
				"name": "libyaml-0-2:amd64",
				"source": {
					"name": "libyaml",
					"upstreamVersion": "0.1.7",
					"version": "0.1.7-2ubuntu3"
				},
				"summary": "Fast YAML 1.1 parser and emitter library",
				"version": "0.1.7-2ubuntu3"
			},
			{
				"arch": "amd64",
				"name": "libzstd1:amd64",
				"source": {
					"name": "libzstd",
					"upstreamVersion": "1.3.3+dfsg",
					"version": "1.3.3+dfsg-2ubuntu1.1"
				},
				"summary": "fast lossless compression algorithm",
				"version": "1.3.3+dfsg-2ubuntu1.1"
			},
			{
				"arch": "amd64",
				"name": "linux-libc-dev:amd64",
				"source": {
					"name": "linux",
					"upstreamVersion": "4.15.0",
					"version": "4.15.0-136.140"
				},
				"summary": "Linux Kernel Headers for development",
				"version": "4.15.0-136.140"
			},
			{
				"arch": "all",
				"name": "locales",
				"source": {
					"name": "glibc",
					"upstreamVersion": "2.27",
					"version": "2.27-3ubuntu1.4"
				},
				"summary": "GNU C Library: National Language (locale) data [support]",
				"version": "2.27-3ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "login",
				"source": {
					"name": "shadow",
					"upstreamVersion": "4.5",
					"version": "1:4.5-1ubuntu2"
				},
				"summary": "system login tools",
				"version": "1:4.5-1ubuntu2"
			},
			{
				"arch": "all",
				"name": "lsb-base",
				"source": {
					"name": "lsb",
					"upstreamVersion": "9.20170808ubuntu1",
					"version": "9.20170808ubuntu1"
				},
				"summary": "Linux Standard Base init script functionality",
				"version": "9.20170808ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "make",
				"source": {
					"name": "make-dfsg",
					"upstreamVersion": "4.1",
					"version": "4.1-9.1ubuntu1"
				},
				"summary": "utility for directing compilation",
				"version": "4.1-9.1ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "mawk",
				"source": {
					"name": "mawk",
					"upstreamVersion": "1.3.3",
					"version": "1.3.3-17ubuntu3"
				},
				"summary": "a pattern scanning and text processing language",
				"version": "1.3.3-17ubuntu3"
			},
			{
				"arch": "amd64",
				"name": "mount",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "tools for mounting and manipulating filesystems",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "all",
				"name": "ncurses-base",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "basic terminal type definitions",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "ncurses-bin",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "terminal-related programs and man pages",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "openssl",
				"source": {
					"name": "openssl",
					"upstreamVersion": "1.1.1",
					"version": "1.1.1-1ubuntu2.1~18.04.8"
				},
				"summary": "Secure Sockets Layer toolkit - cryptographic utility",
				"version": "1.1.1-1ubuntu2.1~18.04.8"
			},
			{
				"arch": "amd64",
				"name": "passwd",
				"source": {
					"name": "shadow",
					"upstreamVersion": "4.5",
					"version": "1:4.5-1ubuntu2"
				},
				"summary": "change and administer password and group data",
				"version": "1:4.5-1ubuntu2"
			},
			{
				"arch": "amd64",
				"name": "patch",
				"source": {
					"name": "patch",
					"upstreamVersion": "2.7.6",
					"version": "2.7.6-2ubuntu1.1"
				},
				"summary": "Apply a diff file to an original",
				"version": "2.7.6-2ubuntu1.1"
			},
			{
				"arch": "amd64",
				"name": "perl",
				"source": {
					"name": "perl",
					"upstreamVersion": "5.26.1",
					"version": "5.26.1-6ubuntu0.5"
				},
				"summary": "Larry Wall's Practical Extraction and Report Language",
				"version": "5.26.1-6ubuntu0.5"
			},
			{
				"arch": "amd64",
				"name": "perl-base",
				"source": {
					"name": "perl",
					"upstreamVersion": "5.26.1",
					"version": "5.26.1-6ubuntu0.5"
				},
				"summary": "minimal Perl system",
				"version": "5.26.1-6ubuntu0.5"
			},
			{
				"arch": "all",
				"name": "perl-modules-5.26",
				"source": {
					"name": "perl",
					"upstreamVersion": "5.26.1",
					"version": "5.26.1-6ubuntu0.5"
				},
				"summary": "Core Perl modules",
				"version": "5.26.1-6ubuntu0.5"
			},
			{
				"arch": "amd64",
				"name": "procps",
				"source": {
					"name": "procps",
					"upstreamVersion": "3.3.12",
					"version": "2:3.3.12-3ubuntu1.2"
				},
				"summary": "/proc file system utilities",
				"version": "2:3.3.12-3ubuntu1.2"
			},
			{
				"arch": "amd64",
				"name": "sed",
				"source": {
					"name": "sed",
					"upstreamVersion": "4.4",
					"version": "4.4-2"
				},
				"summary": "GNU stream editor for filtering/transforming text",
				"version": "4.4-2"
			},
			{
				"arch": "all",
				"name": "sensible-utils",
				"source": {
					"name": "sensible-utils",
					"upstreamVersion": "0.0.12",
					"version": "0.0.12"
				},
				"summary": "Utilities for sensible alternative selection",
				"version": "0.0.12"
			},
			{
				"arch": "amd64",
				"name": "sysvinit-utils",
				"source": {
					"name": "sysvinit",
					"upstreamVersion": "2.88dsf",
					"version": "2.88dsf-59.10ubuntu1"
				},
				"summary": "System-V-like utilities",
				"version": "2.88dsf-59.10ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "tar",
				"source": {
					"name": "tar",
					"upstreamVersion": "1.29b",
					"version": "1.29b-2ubuntu0.2"
				},
				"summary": "GNU version of the tar archiving utility",
				"version": "1.29b-2ubuntu0.2"
			},
			{
				"arch": "all",
				"name": "tzdata",
				"source": {
					"name": "tzdata",
					"upstreamVersion": "2021a",
					"version": "2021a-0ubuntu0.18.04"
				},
				"summary": "time zone and daylight-saving time data",
				"version": "2021a-0ubuntu0.18.04"
			},
			{
				"arch": "all",
				"name": "ubuntu-keyring",
				"source": {
					"name": "ubuntu-keyring",
					"upstreamVersion": "2018.09.18.1~18.04.0",
					"version": "2018.09.18.1~18.04.0"
				},
				"summary": "GnuPG keys of the Ubuntu archive",
				"version": "2018.09.18.1~18.04.0"
			},
			{
				"arch": "amd64",
				"name": "util-linux",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "miscellaneous system utilities",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "xz-utils",
				"source": {
					"name": "xz-utils",
					"upstreamVersion": "5.2.2",
					"version": "5.2.2-1.3"
				},
				"summary": "XZ-format compression utilities",
				"version": "5.2.2-1.3"
			},
			{
				"arch": "amd64",
				"name": "zlib1g:amd64",
				"source": {
					"name": "zlib",
					"upstreamVersion": "1.2.11.dfsg",
					"version": "1:1.2.11.dfsg-0ubuntu2"
				},
				"summary": "compression library - runtime",
				"version": "1:1.2.11.dfsg-0ubuntu2"
			},
			{
				"arch": "amd64",
				"name": "zlib1g-dev:amd64",
				"source": {
					"name": "zlib",
					"upstreamVersion": "1.2.11.dfsg",
					"version": "1:1.2.11.dfsg-0ubuntu2"
				},
				"summary": "compression library - development",
				"version": "1:1.2.11.dfsg-0ubuntu2"
			}
		]`))

		runPackageMetadata, err := bionic.GetRunPackageMetadata(runImageTag)
		Expect(err).NotTo(HaveOccurred())
		Expect(runPackageMetadata).To(MatchJSON(`[
			{
				"arch": "all",
				"name": "adduser",
				"source": {
					"name": "adduser",
					"upstreamVersion": "3.116ubuntu1",
					"version": "3.116ubuntu1"
				},
				"summary": "add and remove users and groups",
				"version": "3.116ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "apt",
				"source": {
					"name": "apt",
					"upstreamVersion": "1.6.12ubuntu0.2",
					"version": "1.6.12ubuntu0.2"
				},
				"summary": "commandline package manager",
				"version": "1.6.12ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "base-files",
				"source": {
					"name": "base-files",
					"upstreamVersion": "10.1ubuntu2.10",
					"version": "10.1ubuntu2.10"
				},
				"summary": "Debian base system miscellaneous files",
				"version": "10.1ubuntu2.10"
			},
			{
				"arch": "amd64",
				"name": "base-passwd",
				"source": {
					"name": "base-passwd",
					"upstreamVersion": "3.5.44",
					"version": "3.5.44"
				},
				"summary": "Debian base system master password and group files",
				"version": "3.5.44"
			},
			{
				"arch": "amd64",
				"name": "bash",
				"source": {
					"name": "bash",
					"upstreamVersion": "4.4.18",
					"version": "4.4.18-2ubuntu1.2"
				},
				"summary": "GNU Bourne Again SHell",
				"version": "4.4.18-2ubuntu1.2"
			},
			{
				"arch": "amd64",
				"name": "bsdutils",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "basic utilities from 4.4BSD-Lite",
				"version": "1:2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "bzip2",
				"source": {
					"name": "bzip2",
					"upstreamVersion": "1.0.6",
					"version": "1.0.6-8.1ubuntu0.2"
				},
				"summary": "high-quality block-sorting file compressor - utilities",
				"version": "1.0.6-8.1ubuntu0.2"
			},
			{
				"arch": "all",
				"name": "ca-certificates",
				"source": {
					"name": "ca-certificates",
					"upstreamVersion": "20210119~18.04.1",
					"version": "20210119~18.04.1"
				},
				"summary": "Common CA certificates",
				"version": "20210119~18.04.1"
			},
			{
				"arch": "amd64",
				"name": "coreutils",
				"source": {
					"name": "coreutils",
					"upstreamVersion": "8.28",
					"version": "8.28-1ubuntu1"
				},
				"summary": "GNU core utilities",
				"version": "8.28-1ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "dash",
				"source": {
					"name": "dash",
					"upstreamVersion": "0.5.8",
					"version": "0.5.8-2.10"
				},
				"summary": "POSIX-compliant shell",
				"version": "0.5.8-2.10"
			},
			{
				"arch": "all",
				"name": "debconf",
				"source": {
					"name": "debconf",
					"upstreamVersion": "1.5.66ubuntu1",
					"version": "1.5.66ubuntu1"
				},
				"summary": "Debian configuration management system",
				"version": "1.5.66ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "debianutils",
				"source": {
					"name": "debianutils",
					"upstreamVersion": "4.8.4",
					"version": "4.8.4"
				},
				"summary": "Miscellaneous utilities specific to Debian",
				"version": "4.8.4"
			},
			{
				"arch": "amd64",
				"name": "diffutils",
				"source": {
					"name": "diffutils",
					"upstreamVersion": "3.6",
					"version": "1:3.6-1"
				},
				"summary": "File comparison utilities",
				"version": "1:3.6-1"
			},
			{
				"arch": "amd64",
				"name": "dpkg",
				"source": {
					"name": "dpkg",
					"upstreamVersion": "1.19.0.5ubuntu2.3",
					"version": "1.19.0.5ubuntu2.3"
				},
				"summary": "Debian package management system",
				"version": "1.19.0.5ubuntu2.3"
			},
			{
				"arch": "amd64",
				"name": "e2fsprogs",
				"source": {
					"name": "e2fsprogs",
					"upstreamVersion": "1.44.1",
					"version": "1.44.1-1ubuntu1.3"
				},
				"summary": "ext2/ext3/ext4 file system utilities",
				"version": "1.44.1-1ubuntu1.3"
			},
			{
				"arch": "amd64",
				"name": "fdisk",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "collection of partitioning utilities",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "findutils",
				"source": {
					"name": "findutils",
					"upstreamVersion": "4.6.0+git+20170828",
					"version": "4.6.0+git+20170828-2"
				},
				"summary": "utilities for finding files--find, xargs",
				"version": "4.6.0+git+20170828-2"
			},
			{
				"arch": "amd64",
				"name": "gcc-8-base:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GCC, the GNU Compiler Collection (base package)",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "gpgv",
				"source": {
					"name": "gnupg2",
					"upstreamVersion": "2.2.4",
					"version": "2.2.4-1ubuntu1.4"
				},
				"summary": "GNU privacy guard - signature verification tool",
				"version": "2.2.4-1ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "grep",
				"source": {
					"name": "grep",
					"upstreamVersion": "3.1",
					"version": "3.1-2build1"
				},
				"summary": "GNU grep, egrep and fgrep",
				"version": "3.1-2build1"
			},
			{
				"arch": "amd64",
				"name": "gzip",
				"source": {
					"name": "gzip",
					"upstreamVersion": "1.6",
					"version": "1.6-5ubuntu1"
				},
				"summary": "GNU compression utilities",
				"version": "1.6-5ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "hostname",
				"source": {
					"name": "hostname",
					"upstreamVersion": "3.20",
					"version": "3.20"
				},
				"summary": "utility to set/show the host name or domain name",
				"version": "3.20"
			},
			{
				"arch": "all",
				"name": "init-system-helpers",
				"source": {
					"name": "init-system-helpers",
					"upstreamVersion": "1.51",
					"version": "1.51"
				},
				"summary": "helper tools for all init systems",
				"version": "1.51"
			},
			{
				"arch": "amd64",
				"name": "libacl1:amd64",
				"source": {
					"name": "acl",
					"upstreamVersion": "2.2.52",
					"version": "2.2.52-3build1"
				},
				"summary": "Access control list shared library",
				"version": "2.2.52-3build1"
			},
			{
				"arch": "amd64",
				"name": "libapt-pkg5.0:amd64",
				"source": {
					"name": "apt",
					"upstreamVersion": "1.6.12ubuntu0.2",
					"version": "1.6.12ubuntu0.2"
				},
				"summary": "package management runtime library",
				"version": "1.6.12ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libattr1:amd64",
				"source": {
					"name": "attr",
					"upstreamVersion": "2.4.47",
					"version": "1:2.4.47-2build1"
				},
				"summary": "Extended attribute shared library",
				"version": "1:2.4.47-2build1"
			},
			{
				"arch": "all",
				"name": "libaudit-common",
				"source": {
					"name": "audit",
					"upstreamVersion": "2.8.2",
					"version": "1:2.8.2-1ubuntu1.1"
				},
				"summary": "Dynamic library for security auditing - common files",
				"version": "1:2.8.2-1ubuntu1.1"
			},
			{
				"arch": "amd64",
				"name": "libaudit1:amd64",
				"source": {
					"name": "audit",
					"upstreamVersion": "2.8.2",
					"version": "1:2.8.2-1ubuntu1.1"
				},
				"summary": "Dynamic library for security auditing",
				"version": "1:2.8.2-1ubuntu1.1"
			},
			{
				"arch": "amd64",
				"name": "libblkid1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "block device ID library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libbz2-1.0:amd64",
				"source": {
					"name": "bzip2",
					"upstreamVersion": "1.0.6",
					"version": "1.0.6-8.1ubuntu0.2"
				},
				"summary": "high-quality block-sorting file compressor library - runtime",
				"version": "1.0.6-8.1ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "libc-bin",
				"source": {
					"name": "glibc",
					"upstreamVersion": "2.27",
					"version": "2.27-3ubuntu1.4"
				},
				"summary": "GNU C Library: Binaries",
				"version": "2.27-3ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "libc6:amd64",
				"source": {
					"name": "glibc",
					"upstreamVersion": "2.27",
					"version": "2.27-3ubuntu1.4"
				},
				"summary": "GNU C Library: Shared libraries",
				"version": "2.27-3ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "libcap-ng0:amd64",
				"source": {
					"name": "libcap-ng",
					"upstreamVersion": "0.7.7",
					"version": "0.7.7-3.1"
				},
				"summary": "An alternate POSIX capabilities library",
				"version": "0.7.7-3.1"
			},
			{
				"arch": "amd64",
				"name": "libcom-err2:amd64",
				"source": {
					"name": "e2fsprogs",
					"upstreamVersion": "1.44.1",
					"version": "1.44.1-1ubuntu1.3"
				},
				"summary": "common error description library",
				"version": "1.44.1-1ubuntu1.3"
			},
			{
				"arch": "amd64",
				"name": "libdb5.3:amd64",
				"source": {
					"name": "db5.3",
					"upstreamVersion": "5.3.28",
					"version": "5.3.28-13.1ubuntu1.1"
				},
				"summary": "Berkeley v5.3 Database Libraries [runtime]",
				"version": "5.3.28-13.1ubuntu1.1"
			},
			{
				"arch": "amd64",
				"name": "libdebconfclient0:amd64",
				"source": {
					"name": "cdebconf",
					"upstreamVersion": "0.213ubuntu1",
					"version": "0.213ubuntu1"
				},
				"summary": "Debian Configuration Management System (C-implementation library)",
				"version": "0.213ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "libext2fs2:amd64",
				"source": {
					"name": "e2fsprogs",
					"upstreamVersion": "1.44.1",
					"version": "1.44.1-1ubuntu1.3"
				},
				"summary": "ext2/ext3/ext4 file system libraries",
				"version": "1.44.1-1ubuntu1.3"
			},
			{
				"arch": "amd64",
				"name": "libfdisk1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "fdisk partitioning library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libffi6:amd64",
				"source": {
					"name": "libffi",
					"upstreamVersion": "3.2.1",
					"version": "3.2.1-8"
				},
				"summary": "Foreign Function Interface library runtime",
				"version": "3.2.1-8"
			},
			{
				"arch": "amd64",
				"name": "libgcc1:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GCC support library",
				"version": "1:8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libgcrypt20:amd64",
				"source": {
					"name": "libgcrypt20",
					"upstreamVersion": "1.8.1",
					"version": "1.8.1-4ubuntu1.2"
				},
				"summary": "LGPL Crypto library - runtime library",
				"version": "1.8.1-4ubuntu1.2"
			},
			{
				"arch": "amd64",
				"name": "libgmp10:amd64",
				"source": {
					"name": "gmp",
					"upstreamVersion": "6.1.2+dfsg",
					"version": "2:6.1.2+dfsg-2"
				},
				"summary": "Multiprecision arithmetic library",
				"version": "2:6.1.2+dfsg-2"
			},
			{
				"arch": "amd64",
				"name": "libgnutls30:amd64",
				"source": {
					"name": "gnutls28",
					"upstreamVersion": "3.5.18",
					"version": "3.5.18-1ubuntu1.4"
				},
				"summary": "GNU TLS library - main runtime library",
				"version": "3.5.18-1ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "libgpg-error0:amd64",
				"source": {
					"name": "libgpg-error",
					"upstreamVersion": "1.27",
					"version": "1.27-6"
				},
				"summary": "library for common error values and messages in GnuPG components",
				"version": "1.27-6"
			},
			{
				"arch": "amd64",
				"name": "libhogweed4:amd64",
				"source": {
					"name": "nettle",
					"upstreamVersion": "3.4",
					"version": "3.4-1"
				},
				"summary": "low level cryptographic library (public-key cryptos)",
				"version": "3.4-1"
			},
			{
				"arch": "amd64",
				"name": "libidn2-0:amd64",
				"source": {
					"name": "libidn2",
					"upstreamVersion": "2.0.4",
					"version": "2.0.4-1.1ubuntu0.2"
				},
				"summary": "Internationalized domain names (IDNA2008/TR46) library",
				"version": "2.0.4-1.1ubuntu0.2"
			},
			{
				"arch": "amd64",
				"name": "liblz4-1:amd64",
				"source": {
					"name": "lz4",
					"upstreamVersion": "0.0~r131",
					"version": "0.0~r131-2ubuntu3"
				},
				"summary": "Fast LZ compression algorithm library - runtime",
				"version": "0.0~r131-2ubuntu3"
			},
			{
				"arch": "amd64",
				"name": "liblzma5:amd64",
				"source": {
					"name": "xz-utils",
					"upstreamVersion": "5.2.2",
					"version": "5.2.2-1.3"
				},
				"summary": "XZ-format compression library",
				"version": "5.2.2-1.3"
			},
			{
				"arch": "amd64",
				"name": "libmount1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "device mounting library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libncurses5:amd64",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "shared libraries for terminal handling",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "libncursesw5:amd64",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "shared libraries for terminal handling (wide character support)",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "libnettle6:amd64",
				"source": {
					"name": "nettle",
					"upstreamVersion": "3.4",
					"version": "3.4-1"
				},
				"summary": "low level cryptographic library (symmetric and one-way cryptos)",
				"version": "3.4-1"
			},
			{
				"arch": "amd64",
				"name": "libp11-kit0:amd64",
				"source": {
					"name": "p11-kit",
					"upstreamVersion": "0.23.9",
					"version": "0.23.9-2ubuntu0.1"
				},
				"summary": "library for loading and coordinating access to PKCS#11 modules - runtime",
				"version": "0.23.9-2ubuntu0.1"
			},
			{
				"arch": "amd64",
				"name": "libpam-modules:amd64",
				"source": {
					"name": "pam",
					"upstreamVersion": "1.1.8",
					"version": "1.1.8-3.6ubuntu2.18.04.2"
				},
				"summary": "Pluggable Authentication Modules for PAM",
				"version": "1.1.8-3.6ubuntu2.18.04.2"
			},
			{
				"arch": "amd64",
				"name": "libpam-modules-bin",
				"source": {
					"name": "pam",
					"upstreamVersion": "1.1.8",
					"version": "1.1.8-3.6ubuntu2.18.04.2"
				},
				"summary": "Pluggable Authentication Modules for PAM - helper binaries",
				"version": "1.1.8-3.6ubuntu2.18.04.2"
			},
			{
				"arch": "all",
				"name": "libpam-runtime",
				"source": {
					"name": "pam",
					"upstreamVersion": "1.1.8",
					"version": "1.1.8-3.6ubuntu2.18.04.2"
				},
				"summary": "Runtime support for the PAM library",
				"version": "1.1.8-3.6ubuntu2.18.04.2"
			},
			{
				"arch": "amd64",
				"name": "libpam0g:amd64",
				"source": {
					"name": "pam",
					"upstreamVersion": "1.1.8",
					"version": "1.1.8-3.6ubuntu2.18.04.2"
				},
				"summary": "Pluggable Authentication Modules library",
				"version": "1.1.8-3.6ubuntu2.18.04.2"
			},
			{
				"arch": "amd64",
				"name": "libpcre3:amd64",
				"source": {
					"name": "pcre3",
					"upstreamVersion": "8.39",
					"version": "2:8.39-9"
				},
				"summary": "Old Perl 5 Compatible Regular Expression Library - runtime files",
				"version": "2:8.39-9"
			},
			{
				"arch": "amd64",
				"name": "libprocps6:amd64",
				"source": {
					"name": "procps",
					"upstreamVersion": "3.3.12",
					"version": "2:3.3.12-3ubuntu1.2"
				},
				"summary": "library for accessing process information from /proc",
				"version": "2:3.3.12-3ubuntu1.2"
			},
			{
				"arch": "amd64",
				"name": "libseccomp2:amd64",
				"source": {
					"name": "libseccomp",
					"upstreamVersion": "2.4.3",
					"version": "2.4.3-1ubuntu3.18.04.3"
				},
				"summary": "high level interface to Linux seccomp filter",
				"version": "2.4.3-1ubuntu3.18.04.3"
			},
			{
				"arch": "amd64",
				"name": "libselinux1:amd64",
				"source": {
					"name": "libselinux",
					"upstreamVersion": "2.7",
					"version": "2.7-2build2"
				},
				"summary": "SELinux runtime shared libraries",
				"version": "2.7-2build2"
			},
			{
				"arch": "all",
				"name": "libsemanage-common",
				"source": {
					"name": "libsemanage",
					"upstreamVersion": "2.7",
					"version": "2.7-2build2"
				},
				"summary": "Common files for SELinux policy management libraries",
				"version": "2.7-2build2"
			},
			{
				"arch": "amd64",
				"name": "libsemanage1:amd64",
				"source": {
					"name": "libsemanage",
					"upstreamVersion": "2.7",
					"version": "2.7-2build2"
				},
				"summary": "SELinux policy management library",
				"version": "2.7-2build2"
			},
			{
				"arch": "amd64",
				"name": "libsepol1:amd64",
				"source": {
					"name": "libsepol",
					"upstreamVersion": "2.7",
					"version": "2.7-1"
				},
				"summary": "SELinux library for manipulating binary security policies",
				"version": "2.7-1"
			},
			{
				"arch": "amd64",
				"name": "libsmartcols1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "smart column output alignment library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libss2:amd64",
				"source": {
					"name": "e2fsprogs",
					"upstreamVersion": "1.44.1",
					"version": "1.44.1-1ubuntu1.3"
				},
				"summary": "command-line interface parsing library",
				"version": "1.44.1-1ubuntu1.3"
			},
			{
				"arch": "amd64",
				"name": "libssl1.1:amd64",
				"source": {
					"name": "openssl",
					"upstreamVersion": "1.1.1",
					"version": "1.1.1-1ubuntu2.1~18.04.8"
				},
				"summary": "Secure Sockets Layer toolkit - shared libraries",
				"version": "1.1.1-1ubuntu2.1~18.04.8"
			},
			{
				"arch": "amd64",
				"name": "libstdc++6:amd64",
				"source": {
					"name": "gcc-8",
					"upstreamVersion": "8.4.0",
					"version": "8.4.0-1ubuntu1~18.04"
				},
				"summary": "GNU Standard C++ Library v3",
				"version": "8.4.0-1ubuntu1~18.04"
			},
			{
				"arch": "amd64",
				"name": "libsystemd0:amd64",
				"source": {
					"name": "systemd",
					"upstreamVersion": "237",
					"version": "237-3ubuntu10.44"
				},
				"summary": "systemd utility library",
				"version": "237-3ubuntu10.44"
			},
			{
				"arch": "amd64",
				"name": "libtasn1-6:amd64",
				"source": {
					"name": "libtasn1-6",
					"upstreamVersion": "4.13",
					"version": "4.13-2"
				},
				"summary": "Manage ASN.1 structures (runtime)",
				"version": "4.13-2"
			},
			{
				"arch": "amd64",
				"name": "libtinfo5:amd64",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "shared low-level terminfo library for terminal handling",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "libudev1:amd64",
				"source": {
					"name": "systemd",
					"upstreamVersion": "237",
					"version": "237-3ubuntu10.44"
				},
				"summary": "libudev shared library",
				"version": "237-3ubuntu10.44"
			},
			{
				"arch": "amd64",
				"name": "libunistring2:amd64",
				"source": {
					"name": "libunistring",
					"upstreamVersion": "0.9.9",
					"version": "0.9.9-0ubuntu2"
				},
				"summary": "Unicode string library for C",
				"version": "0.9.9-0ubuntu2"
			},
			{
				"arch": "amd64",
				"name": "libuuid1:amd64",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "Universally Unique ID library",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "libyaml-0-2:amd64",
				"source": {
					"name": "libyaml",
					"upstreamVersion": "0.1.7",
					"version": "0.1.7-2ubuntu3"
				},
				"summary": "Fast YAML 1.1 parser and emitter library",
				"version": "0.1.7-2ubuntu3"
			},
			{
				"arch": "amd64",
				"name": "libzstd1:amd64",
				"source": {
					"name": "libzstd",
					"upstreamVersion": "1.3.3+dfsg",
					"version": "1.3.3+dfsg-2ubuntu1.1"
				},
				"summary": "fast lossless compression algorithm",
				"version": "1.3.3+dfsg-2ubuntu1.1"
			},
			{
				"arch": "all",
				"name": "locales",
				"source": {
					"name": "glibc",
					"upstreamVersion": "2.27",
					"version": "2.27-3ubuntu1.4"
				},
				"summary": "GNU C Library: National Language (locale) data [support]",
				"version": "2.27-3ubuntu1.4"
			},
			{
				"arch": "amd64",
				"name": "login",
				"source": {
					"name": "shadow",
					"upstreamVersion": "4.5",
					"version": "1:4.5-1ubuntu2"
				},
				"summary": "system login tools",
				"version": "1:4.5-1ubuntu2"
			},
			{
				"arch": "all",
				"name": "lsb-base",
				"source": {
					"name": "lsb",
					"upstreamVersion": "9.20170808ubuntu1",
					"version": "9.20170808ubuntu1"
				},
				"summary": "Linux Standard Base init script functionality",
				"version": "9.20170808ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "mawk",
				"source": {
					"name": "mawk",
					"upstreamVersion": "1.3.3",
					"version": "1.3.3-17ubuntu3"
				},
				"summary": "a pattern scanning and text processing language",
				"version": "1.3.3-17ubuntu3"
			},
			{
				"arch": "amd64",
				"name": "mount",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "tools for mounting and manipulating filesystems",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "all",
				"name": "ncurses-base",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "basic terminal type definitions",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "ncurses-bin",
				"source": {
					"name": "ncurses",
					"upstreamVersion": "6.1",
					"version": "6.1-1ubuntu1.18.04"
				},
				"summary": "terminal-related programs and man pages",
				"version": "6.1-1ubuntu1.18.04"
			},
			{
				"arch": "amd64",
				"name": "openssl",
				"source": {
					"name": "openssl",
					"upstreamVersion": "1.1.1",
					"version": "1.1.1-1ubuntu2.1~18.04.8"
				},
				"summary": "Secure Sockets Layer toolkit - cryptographic utility",
				"version": "1.1.1-1ubuntu2.1~18.04.8"
			},
			{
				"arch": "amd64",
				"name": "passwd",
				"source": {
					"name": "shadow",
					"upstreamVersion": "4.5",
					"version": "1:4.5-1ubuntu2"
				},
				"summary": "change and administer password and group data",
				"version": "1:4.5-1ubuntu2"
			},
			{
				"arch": "amd64",
				"name": "perl-base",
				"source": {
					"name": "perl",
					"upstreamVersion": "5.26.1",
					"version": "5.26.1-6ubuntu0.5"
				},
				"summary": "minimal Perl system",
				"version": "5.26.1-6ubuntu0.5"
			},
			{
				"arch": "amd64",
				"name": "procps",
				"source": {
					"name": "procps",
					"upstreamVersion": "3.3.12",
					"version": "2:3.3.12-3ubuntu1.2"
				},
				"summary": "/proc file system utilities",
				"version": "2:3.3.12-3ubuntu1.2"
			},
			{
				"arch": "amd64",
				"name": "sed",
				"source": {
					"name": "sed",
					"upstreamVersion": "4.4",
					"version": "4.4-2"
				},
				"summary": "GNU stream editor for filtering/transforming text",
				"version": "4.4-2"
			},
			{
				"arch": "all",
				"name": "sensible-utils",
				"source": {
					"name": "sensible-utils",
					"upstreamVersion": "0.0.12",
					"version": "0.0.12"
				},
				"summary": "Utilities for sensible alternative selection",
				"version": "0.0.12"
			},
			{
				"arch": "amd64",
				"name": "sysvinit-utils",
				"source": {
					"name": "sysvinit",
					"upstreamVersion": "2.88dsf",
					"version": "2.88dsf-59.10ubuntu1"
				},
				"summary": "System-V-like utilities",
				"version": "2.88dsf-59.10ubuntu1"
			},
			{
				"arch": "amd64",
				"name": "tar",
				"source": {
					"name": "tar",
					"upstreamVersion": "1.29b",
					"version": "1.29b-2ubuntu0.2"
				},
				"summary": "GNU version of the tar archiving utility",
				"version": "1.29b-2ubuntu0.2"
			},
			{
				"arch": "all",
				"name": "tzdata",
				"source": {
					"name": "tzdata",
					"upstreamVersion": "2021a",
					"version": "2021a-0ubuntu0.18.04"
				},
				"summary": "time zone and daylight-saving time data",
				"version": "2021a-0ubuntu0.18.04"
			},
			{
				"arch": "all",
				"name": "ubuntu-keyring",
				"source": {
					"name": "ubuntu-keyring",
					"upstreamVersion": "2018.09.18.1~18.04.0",
					"version": "2018.09.18.1~18.04.0"
				},
				"summary": "GnuPG keys of the Ubuntu archive",
				"version": "2018.09.18.1~18.04.0"
			},
			{
				"arch": "amd64",
				"name": "util-linux",
				"source": {
					"name": "util-linux",
					"upstreamVersion": "2.31.1",
					"version": "2.31.1-0.4ubuntu3.7"
				},
				"summary": "miscellaneous system utilities",
				"version": "2.31.1-0.4ubuntu3.7"
			},
			{
				"arch": "amd64",
				"name": "zlib1g:amd64",
				"source": {
					"name": "zlib",
					"upstreamVersion": "1.2.11.dfsg",
					"version": "1:1.2.11.dfsg-0ubuntu2"
				},
				"summary": "compression library - runtime",
				"version": "1:1.2.11.dfsg-0ubuntu2"
			}
		]`))
	})
}
