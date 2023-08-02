import socket
import random
import time
import sys


def sendTCP():
    # Define the IP address and port of the destination server.
    ip_address = "sender"
    port = 8000
    # Create a socket object.
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Connect to the destination server.
    sock.connect((ip_address, port))
    print(f"Addr = {ip_address}:{port}")
    # Define the data to be sent.
    data = f""
    for i in range(5):
        data += (f"id{str(i).zfill(2)}")
        data += (f"type{str(2).zfill(2)}")
        number = round(random.uniform(10, 20), 2)
        if len(str(number)) < 5:
            data += (f"data{str(number).zfill(4)}" + "0")
            # print(f"data{str(number).zfill(4)}" + "0")
        else:
            data += (f"data{str(number).zfill(4)}")
            # print(f"data{str(number).zfill(4)}")

    # Send the data.
    data = f"{len(data)}".zfill(4) + data
    bytes = data.encode(sys.getdefaultencoding())
    sock.sendall(bytes)
    print(f"Msg Send: {data}\t{bytes}")

    # Close the socket.
    sock.close()


if __name__ == "__main__":
    for i in range(100):
        sendTCP()
        time.sleep(0.0001)
