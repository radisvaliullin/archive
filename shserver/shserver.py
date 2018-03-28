#
import cgi
import hashlib
from http.server import BaseHTTPRequestHandler, HTTPServer


class ShortingServer:

    def __init__(self, host, port):
        self._host = host
        self._port = port

        self._hash_url_store = {}

    def run(self):

        server_address = (self._host, self._port)
        server = HTTPServer(server_address, self._make_handler_class())
        server.serve_forever()

    def _make_handler_class(self):
        return make_shserver_request_handler_class(self._host, self._port, self._hash_url_store)


def make_shserver_request_handler_class(host, port, store):

    class ShServerRequestHandler(BaseHTTPRequestHandler):

        _host = host
        _port = port
        _hash_url_store = store

        def do_GET(self):

            # main page
            if self.path == '/':
                self.send_response(200)
                self.send_header('Content-type', 'text/html')
                self.end_headers()

                # Send the html message
                message = self.get_main_page_html()
                self.wfile.write(bytes(message, "utf8"))
                return

            else:

                # redirect by short link
                if "?" not in self.path and len(self.path.split("/")) == 2:
                    hash = self.path.split("/")[1]

                    redirect_url = ""

                    # if short link exist in store
                    if hash in self._hash_url_store:
                        original_url = self._hash_url_store[hash]
                        if not original_url.startswith("http"):
                            original_url = "http://" + original_url
                        redirect_url = original_url
                    else:
                        redirect_url = "http://" + self._host + ":" + str(self._port)

                    self.send_response(200)
                    self.send_header('Content-type', 'text/html')
                    self.end_headers()

                    # Send the html message
                    message = self.get_redirect_html_template().format(redirect_url)
                    self.wfile.write(bytes(message, "utf8"))
                    return

                # wrong link, return main page
                else:
                    self.send_response(200)
                    self.send_header('Content-type', 'text/html')
                    self.end_headers()

                    # Send the html message
                    _url = "http://" + self._host + ":" + str(self._port)
                    message = self.get_redirect_html_template().format(_url)
                    self.wfile.write(bytes(message, "utf8"))
                    return

        def do_POST(self):

            # Parse the form data posted
            form = cgi.FieldStorage(
                fp=self.rfile,
                headers=self.headers,
                environ={
                    'REQUEST_METHOD': 'POST',
                    'CONTENT_TYPE': self.headers['Content-Type'],
                }
            )

            self.send_response(200)
            self.send_header('Content-type', 'text/html')
            self.end_headers()

            # Send the html message
            original_url = ""
            shorted_url = ""
            if "original_url" in form:
                original_url = form["original_url"].value

                url_hasg = hashlib.md5(original_url.encode())
                url_hasg_hex = url_hasg.hexdigest()
                self._hash_url_store[url_hasg_hex] = original_url
                shorted_url = self._host + ":" + str(self._port) + "/" + url_hasg_hex

            message = self.get_shorted_url_html_template().format(shorted_url)

            self.wfile.write(bytes(message, "utf8"))
            return

        def get_main_page_html(self):
            return '''
                <!DOCTYPE HTML>
                <html>
                 <head>
                  <meta charset="utf-8">
                  <title>Short URL</title>
                 </head>
                 <body>
                
                 <form method="post">
                  <p><b>Введите URL для получения короткой ссылки</b></p>
                  <p><input type="text" name="original_url"> <input type="submit"></p>
                 </form>
                
                 </body>
                </html>
            '''

        def get_shorted_url_html_template(self):
            return '''
                <!DOCTYPE HTML>
                <html>
                 <head>
                  <meta charset="utf-8">
                  <title>Short URL</title>
                 </head>
                 <body>
                
                 <form method="post">
                  <p><b>Введите URL для получения короткой ссылки</b></p>
                  <p><input type="text" name="original_url"> <input type="submit"></p>
                  <p><b>Ваш короткий URL - {0}</b></p>
                 </form>
                
                 </body>
                </html>
            '''

        def get_redirect_html_template(self):
            return '''
            <!DOCTYPE HTML>
            <html>
             <head>
              <meta charset="utf-8">
              <meta http-equiv="refresh" content="0; url={0}" />
             </head>
            </html>
            '''

    return ShServerRequestHandler
