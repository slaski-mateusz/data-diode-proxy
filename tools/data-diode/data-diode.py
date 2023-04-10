from os import environ
from ipaddress import ip_address
from socket import socket, AF_INET, SOCK_DGRAM


env_vars_issues = []


try:
    datadiode_in_ip =environ['DATADIODE_IN_IP']
    if datadiode_in_ip == "":
        env_vars_issues.append("DATADIODE_IN_IP is empty")
    else:
        try:
            ip_address(datadiode_in_ip)
        except ValueError:
            env_vars_issues.append(f"Environment variable DATADIODE_IN_IP value {datadiode_in_ip} is not proper ip address")
except KeyError:
    env_vars_issues.append('Environment variable DATADIODE_IN_IP is not defined')

try:
    datadiode_in_port =environ['DATADIODE_IN_PORT']
    if datadiode_in_port == "":
        env_vars_issues.append("DATADIODE_IN_PORT is empty")
    else:
        try:
            datadiode_in_port = int(datadiode_in_port)
        except ValueError:
            env_vars_issues.append(f"DATADIODE_IN_PORT value {datadiode_in_port} is not integer number")
except KeyError:
    env_vars_issues.append('Environment variable DATADIODE_IN_PORT is not defined')


try:
    datadiode_out_ip =environ['DATADIODE_OUT_IP']
    if datadiode_out_ip == "":
        env_vars_issues.append("DATADIODE_OUT_IP is empty")
    else:
        try:
            ip_address(datadiode_out_ip)
        except ValueError:
            env_vars_issues.append(f"Environment variable DATADIODE_OUT_IP value {datadiode_out_ip} is not proper ip address")
except KeyError:
    env_vars_issues.append('Environment variable DATADIODE_OUT_IP is not defined')

try:
    datadiode_out_port =environ['DATADIODE_OUT_PORT']
    if datadiode_out_port == "":
        env_vars_issues.append("DATADIODE_OUT_PORT is empty")
    else:
        try:
            datadiode_out_port = int(datadiode_out_port)
        except ValueError:
            env_vars_issues.append(f"DATADIODE_OUT_PORT value {datadiode_out_port} is not integer number")
except KeyError:
    env_vars_issues.append('Environment variable DATADIODE_OUT_PORT is not defined')


if env_vars_issues:
    print("Environment variables not defined properly!")
    print("\n".join(env_vars_issues))
    print("Exiting!")
    exit(1)

print(f"Listening for UDP packets on {datadiode_in_ip}:{datadiode_in_port}")
# print(f"Forwarding UDP packets to {datadiode_out_ip}:{datadiode_out_port}")


def main():
    sk = socket(AF_INET, SOCK_DGRAM)
    sk.bind(
        (datadiode_in_ip, datadiode_in_port)
    )
    while True:
        data, addr = sk.recvfrom(1024)
        print("from {} received data {}".format(addr, data))
        

if __name__ == "__main__":
    main()