from kubernetes import client, config
import base64, argparse

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Extract, in a cross-platform way, secrets from kubernetes")
    parser.add_argument('--namespace', dest='namespace', default='default')
    parser.add_argument('--name', dest='name', default='rabbitmq')
    parser.add_argument('--secret', dest='secret', help="For the default rabbitmq extraction, secrets are 'rabbitmq-password' and 'rabbitmq-erlang-cookie'")

    args = parser.parse_args()

    config.load_kube_config()
    v1 = client.CoreV1Api()
    sec = v1.read_namespaced_secret(args.name, args.namespace)
    encoded = sec.data[args.secret]
    print(base64.b64decode(encoded).decode('ascii'))
