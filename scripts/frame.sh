
# update apt and install libx11-dev and anything needed for chromium to 
 apt update && apt install -y libx11-dev
 

#  If the clipboard package is in an environment without a frame buffer,
#  such as a cloud server, it may also be necessary to install xvfb:
 
 apt install -y xvfb
 
#  and initialize a virtual frame buffer:
 
 Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
 export DISPLAY=:99.0