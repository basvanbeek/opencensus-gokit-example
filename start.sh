#!/bin/sh
#nohup build/ocg-elegantmonolith &>elegantmonolith.log &
nohup build/ocg-qrgenerator     &>qrgenerator.log     &
nohup build/ocg-device          &>device.log          &
nohup build/ocg-event           &>event.log           &
nohup build/ocg-frontend        &>frontend.log        &
