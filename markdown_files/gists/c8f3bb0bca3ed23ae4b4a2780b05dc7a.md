# Setup DNSMASQ

[View original Gist on GitHub](https://gist.github.com/Integralist/c8f3bb0bca3ed23ae4b4a2780b05dc7a)

## Setup DNSMASQ.sh

```shell
function configure_dnsmasq {
  # https://help.ubuntu.com/community/Dnsmasq
  apt-get install -y dnsmasq
  sed -i 's/#listen-address=/listen-address=127.0.0.1/' /etc/dnsmasq.conf
  sed -i 's/#user=/user=root/' /etc/dnsmasq.conf
  sed -i 's/#prepend domain-name-servers 127.0.0.1;/prepend domain-name-servers 127.0.0.1;/' /etc/dhcp/dhclient.conf

  # You can't use sed -i on a file that's bind mounted by Docker
  # But you can 'edit' it still...
  vi -E -s /etc/resolv.conf <<-EOF
  :1s/^/nameserver 127.0.0.1\r/
  :update
  :quit
EOF

  echo 'addn-hosts=/etc/hosts' >> /etc/dnsmasq.d/hosts.conf
  /etc/init.d/dnsmasq restart
}
```

## Update etc hosts.py

```python
import yaml


def upstreams(path):
    with open(path, 'r') as f:
        config = yaml.load(f)['default']['config']
        return config['upstreams']


def parse_host(v):
    return v.split(':')[0]


def sort(ups):
    s = set()
    for k, v in ups.items():
        if isinstance(v, list):
            for i in v:
                s.add(parse_host(i))
        else:
            s.add(parse_host(v))
    return s


def update_hosts(dns):
    with open('/etc/hosts', 'a') as f:
        f.write('\n######################## new hosts added by test suite #########################\n\n')
        for host in dns:
            f.write('127.0.0.1 {}\n'.format(host))


def show_hosts():
    print('================================= /etc/hosts ===================================')
    with open('/etc/hosts', 'r') as f:
        print(f.read())

update_hosts(sort(upstreams('app/config.yml')))
show_hosts()
```

