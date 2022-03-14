from datetime import datetime

if __name__ == "__main__":
    print(datetime.now().astimezone().replace(microsecond=0).isoformat())
