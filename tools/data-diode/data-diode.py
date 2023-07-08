from os import environ
from ipaddress import ip_address
from socket import socket, AF_INET, SOCK_DGRAM
from random import random

env_vars ={
    "DATADIODE_IN_IP": "IP",
    "DATADIODE_IN_PORT": "PORT",
    "DATADIODE_OUT_IP": "IP",
    "DATADIODE_OUT_PORT": "PORT",
    "DATADIODE_DEST_IP": "IP",
    "DATADIODE_DEST_PORT": "PORT",
    "DATADIODE_DROP_PCT": "int",
    "DATADIODE_STOP_CYCLE": "int",
    "DATADIODE_STOP_RND": "int",
    "DATADIODE_STOP_MIN": "int",
    "DATADIODE_STOP_MAX": "int"
}

env_vars_issues = []

def get_env_value(env_var_name, var_content, env_vars_issues):
    try:
        env_var_value =environ[env_var_name]
        if env_var_value == "":
            env_vars_issues.append(f"Environment variable {env_var_name} is empty")
        else:
            if var_content == "IP":
                try:
                    ip_address(env_var_value)
                    return env_var_value
                except ValueError:
                    env_vars_issues.append(f"Environment variable {env_var_name} value {env_var_value} is not proper ip address")
            elif var_content == "PORT":
                try:
                    port_num = int(env_var_value)
                    if port_num < 1000 or port_num > 65535:
                        raise TypeError
                    return port_num
                except ValueError:
                    env_vars_issues.append(f"Environment variable {env_var_name} value {env_var_value} is not proper port number")
                except TypeError:
                    env_vars_issues.append(f"Environment variable {env_var_name} value {env_var_value} is not proper port range 1000~65535")
            elif var_content == "int":
                try:
                    var_int_value = int(env_var_value)
                    return var_int_value
                except ValueError:
                    env_vars_issues.append(f"Environment variable {env_var_name} value {env_var_value} is not integer")
            else:
                env_vars_issues.append(f"Checking environment variables types: IP,PORT,int not {var_content}")
    except KeyError:
        env_vars_issues.append(f"Environment variable {env_var_name} is not defined")


def main():

    datadiode_in_ip = get_env_value(
        "DATADIODE_IN_IP",
        env_vars["DATADIODE_IN_IP"],
        env_vars_issues
    )
    datadiode_in_port = get_env_value(
        "DATADIODE_IN_PORT",
        env_vars["DATADIODE_IN_PORT"],
        env_vars_issues
    )
    datadiode_out_ip = get_env_value(
        "DATADIODE_OUT_IP",
        env_vars["DATADIODE_OUT_IP"],
        env_vars_issues
    )
    datadiode_out_port = get_env_value(
        "DATADIODE_OUT_PORT",
        env_vars["DATADIODE_OUT_PORT"],
        env_vars_issues
    )

    datadiode_dest_ip = get_env_value(
        "DATADIODE_DEST_IP",
        env_vars["DATADIODE_DEST_IP"],
        env_vars_issues
    )
    datadiode_dest_port = get_env_value(
        "DATADIODE_DEST_PORT",
        env_vars["DATADIODE_DEST_PORT"],
        env_vars_issues
    )

    datadiode_drop_pct = get_env_value(
        "DATADIODE_DROP_PCT",
        env_vars["DATADIODE_DROP_PCT"],
        env_vars_issues
    )
    datadiode_stop_cycle = get_env_value(
        "DATADIODE_STOP_CYCLE",
        env_vars["DATADIODE_STOP_CYCLE"],
        env_vars_issues
    )
    datadiode_stop_rnd = get_env_value(
        "DATADIODE_STOP_RND",
        env_vars["DATADIODE_STOP_RND"],
        env_vars_issues
    )
    datadiode_stop_min = get_env_value(
        "DATADIODE_STOP_MIN",
        env_vars["DATADIODE_STOP_MIN"],
        env_vars_issues
    )
    datadiode_stop_max = get_env_value(
        "DATADIODE_STOP_MAX",
        env_vars["DATADIODE_STOP_MAX"],
        env_vars_issues
    )
    

    if env_vars_issues:
        print("Environment variables not defined properly!")
        print("\n".join(env_vars_issues))
        print("Exiting!")
        exit(1)

    print(f"Listening for UDP packets on {datadiode_in_ip}:{datadiode_in_port}")
    # print(f"Forwarding UDP packets to {datadiode_out_ip}:{datadiode_out_port}")

    sk_in = socket(AF_INET, SOCK_DGRAM)
    sk_in.bind(
        (datadiode_in_ip, datadiode_in_port)
    )

    sk_out = socket(AF_INET, SOCK_DGRAM)
    sk_out.bind(
        (datadiode_out_ip, datadiode_out_port)
    )

    while True:
        data, addr = sk_in.recvfrom(1024)
        print("from {} received data {}".format(addr, data))
        drop_packets = False
        if datadiode_drop_pct > 0:
            dpc = 100 * random()
            if dpc < datadiode_drop_pct:
                drop_packets = True
        if not drop_packets:            
            sk_out.sendto(data, ("127.0.0.1", 9000))
        else:
            print("dropping packets")

if __name__ == "__main__":
    main()