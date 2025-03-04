from telemetry.udp_listener import UDPListener

udp = UDPListener()

def main():
    udp.startup_listener()

if __name__ == '__main__':
    main()
