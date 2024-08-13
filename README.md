## AutoEMU ##
A terminal based tool built in Golang as part of my Part III Individual Project.  
Completed while studying for the degree of MEng Computer Science at the University of Southampton.

Thank you to my supervisor, Dr Nawfal Fadhel, for his continued encouragement and support for the duration of this project.

The project works to bring together all the individual tools required for emulating a device into a single, automated application.
This hope to simplify the process to encourage device emulation and virtualisation.  
The tool is far from complete and work will be continued on it.

Currently the tool is only compatiable with Debian/Ubuntu Linux distributions. 

## Introduction ##

This project is working to automate device emulation with no user interaction.

## Prerequisites ##

The following tools should be installed:
```sh
sudo apt-get install git build-essential zlib1g-dev liblzma-dev python3-magic autoconf python-is-python3 xxd binutils lxterminal
```

For QEMU:
```sh
sudo apt-get install qemu-system qemu-user-static
```

For BuildRoot
```sh
sudo apt-get install libncurses-dev bison flex libssl-dev libelf-dev make
```

## Tools Used ##

The tool uses the firmware-mod-kit to extract the root file system from a firmware binary.
The tool can be found at https://github.com/rampageX/firmware-mod-kit
The extract-firmware.sh script has been slightly modified to fit the project needs

The readelf tool is used to analyse ELF executables found within the firmware, to extract the architecture that the firmware was built for.
This hopes to avoid any required user interaction to virtualise a device.

The QEMU tool is used to provide full system emulation.
The tool is available at https://www.qemu.org/

The BuildRoot tool is used to build custom kernels for the emulated devices.
It is packaged within the WorkSpace of AutoEMU but is also available for download at https://buildroot.org/download.html


## Using the Tool ##

There are two ways to use this tool.
It can be used as an interactive terminal tool or as a single command with tags.

To Use the tool as an interactive terminal, run the AutoEMU.exe binary.

## Installing ##

Download the latest zip file in the releases section (not currently avaliable).

This contains the AutoEMU binary and a WorkSpace directory, containing a directory tree.

In this directory, the contents of the firmware-mod-kit repository are also mentioned above, with a modified extract-firmware.sh script.
Once downloaded, run the binary, and it will build a configuration.

When asked for the path to the workspace directory, ensure /WorkSpace is included in the path unless the /WorkSpace directory is present in your current working directory.

The easiest way to install AutoEMU is to download the zip file in the releases section (not currently available), unzip the file and run
```sh
./AutoEMU -config -pwd -pwd
```

## Known Issues ##

Currently, AutoEMU will not determine the machine board which should be used in the QEMU virtual machine.  

## Planned Updates ##
- Addition of a GUI
- Improved extraction of the root file system, kernel and initramfs
- Support for a wider range of filesystems
- Decreased dependence on 3rd party tools, in particular FMK and Buildroot
