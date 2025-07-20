import socket

UniverseCookie = "20050119"

class SocketConn:
    conn=None
    chunk=b''
    def read(self):
        while True:
            pos=self.chunk.find(b'\n')
            if pos!=-1:
                line=self.chunk[:pos]
                self.chunk=self.chunk[pos+1:]
                ret=line.decode('utf-8')
                ret=ret.replace(' ','')
                if ret:
                    return ret
            else:
                data=self.conn.recv(4096)
                if not data:
                    break
                self.chunk+=data
    def write(self,s):
        self.conn.sendall((s+"\n").encode("utf-8"))
    def __init__(self,port):
        self.conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server_address = ('localhost', port)
        self.conn.connect(server_address)