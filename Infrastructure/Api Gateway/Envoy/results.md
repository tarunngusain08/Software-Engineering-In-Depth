Test load balancing:
```sh
for i in {1..10}; do curl -H "Host: www.example.com" http://127.0.0.1:80/; echo ""; done
```

<img width="374" alt="Screenshot 2025-03-17 at 6 32 57 PM" src="https://github.com/user-attachments/assets/a3b4153c-2db2-4516-ae69-b323b9dcec17" />
<img width="894" alt="Screenshot 2025-03-17 at 6 32 03 PM" src="https://github.com/user-attachments/assets/5d28f34d-9d49-42af-8980-424d309156ba" />
<img width="922" alt="Screenshot 2025-03-17 at 6 07 56 PM" src="https://github.com/user-attachments/assets/1b375bf6-0674-4b26-b055-a41a2579823d" />
<img width="735" alt="Screenshot 2025-03-17 at 6 07 44 PM" src="https://github.com/user-attachments/assets/d6302a3b-170c-4cf3-88f9-e15ef2f1daf9" />
<img width="757" alt="Screenshot 2025-03-17 at 6 07 24 PM" src="https://github.com/user-attachments/assets/701811ef-c4d7-4db7-97b1-a25683fb7c0e" />
<img width="1507" alt="Screenshot 2025-03-17 at 6 07 15 PM" src="https://github.com/user-attachments/assets/5e64013d-99c1-4486-9552-8ed9d0bc7548" />
