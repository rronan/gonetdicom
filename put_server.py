import http.server
import sys


class HTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def do_PUT(self):
        path = self.translate_path(self.path)
        length = int(self.headers["Content-Length"])
        with open(path, "wb") as f:
            f.write(self.rfile.read(length))
        self.send_response(200, "Created")
        self.end_headers()


if __name__ == "__main__":
    http.server.test(HandlerClass=HTTPRequestHandler, port=sys.argv[1])
