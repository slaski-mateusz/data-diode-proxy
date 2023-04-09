from socket import socket, AF_INET, SOCK_DGRAM

rec_udp_ip = "127.0.0.1"
rec_udp_port = 2345

def main():
    sk = socket(AF_INET, SOCK_DGRAM)
    sk.bind(
        (rec_udp_ip, rec_udp_port)
    )
    while True:
        data, addr = sk.recvfrom(1024)
        print("from {} received data {}".format(addr, data))
        

if __name__ == "__main__":
    main()