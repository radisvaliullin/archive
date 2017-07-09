from shserver.shserver import ShortingServer

if __name__ == "__main__":

    server = ShortingServer("0.0.0.0", 8000)
    server.run()
