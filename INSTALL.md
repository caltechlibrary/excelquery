
# Installation

*xlquery* is a command line program run from a shell like Bash. You can find compiled
version in the [releases](https://github.com/caltechlibrary/xlquery/releases/latest) 
in the Github repository in a zip file named *xlquery-binary-release.zip*. Inside
the zip file look for the directory that matches your computer and copy that someplace
defined in your path (e.g. $HOME/bin). 

Compiled versions are available for Mac OS X (amd64 processor), Linux (amd64), Windows
(amd64) and Rapsberry Pi (both ARM6 and ARM7)

## Mac OS X

1. Go to [github.com/caltechlibrary/xlquery/releases/latest](https://github.com/caltechlibrary/xlquery/releases/latest)
2. Click on the green "xlquery-binary-release.zip" link and download
3. Open a finder window and find the downloaded file and unzip it (e.g. xlquery-binary-release.zip)
4. Look in the unziped folder and find dist/macosx-amd64/xlquery
5. Drag (or copy) the *xlquery* to a "bin" directory in your path
6. Open and "Terminal" and run `xlquery -h`

## Windows

1. Go to [github.com/caltechlibrary/xlquery/releases/latest](https://github.com/caltechlibrary/xlquery/releases/latest)
2. Click on the green "xlquery-binary-release.zip" link and download
3. Open the file manager find the downloaded file and unzip it (e.g. xlquery-binary-release.zip)
4. Look in the unziped folder and find dist/windows-amd64/xlquery.exe
5. Drag (or copy) the *xlquery.exe* to a "bin" directory in your path
6. Open Bash and and run `xlquery -h`

## Linux

1. Go to [github.com/caltechlibrary/xlquery/releases/latest](https://github.com/caltechlibrary/xlquery/releases/latest)
2. Click on the green "xlquery-binary-release.zip" link and download
3. find the downloaded zip file and unzip it (e.g. unzip ~/Downloads/xlquery-binary-release.zip)
4. In the unziped directory and find for dist/linux-amd64/xlquery
5. copy the *xlquery* to a "bin" directory (e.g. cp ~/Downloads/xlquery-binary-release/dist/linux-amd64/xlquery ~/bin/)
6. From the shell prompt run `xlquery -h`

## Raspberry Pi

If you are using a Raspberry Pi 2 or later use the ARM7 binary, ARM6 is only for the first generaiton Raspberry Pi.

1. Go to [github.com/caltechlibrary/xlquery/releases/latest](https://github.com/caltechlibrary/xlquery/releases/latest)
2. Click on the green "xlquery-binary-release.zip" link and download
3. find the downloaded zip file and unzip it (e.g. unzip ~/Downloads/xlquery-binary-release.zip)
4. In the unziped directory and find for dist/raspberrypi-arm7/xlquery
5. copy the *xlquery* to a "bin" directory (e.g. cp ~/Downloads/xlquery-binary-release/dist/raspberrypi-arm7/xlquery ~/bin/)
    + if you are using an original Raspberry Pi you should copy the ARM6 version instead
6. From the shell prompt run `xlquery -h`

