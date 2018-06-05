#!/bin/sh
cd build
nohup ./ocg-qrgenerator     &>qrgenerator.log     &
nohup ./ocg-device          &>device.log          &
nohup ./ocg-event           &>event.log           &
nohup ./ocg-frontend        &>frontend.log        &
