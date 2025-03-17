Test load balancing:
```sh
for i in {1..10}; do curl -H "Host: www.example.com" http://127.0.0.1:80/; echo ""; done
```

