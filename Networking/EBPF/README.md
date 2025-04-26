# eBPF Active Connections Monitor

This project demonstrates how to use an eBPF program written in Go to monitor active network connections, displaying their PIDs, ports, and other relevant information. Since macOS does not natively support eBPF, we use Lima to create a Linux VM for running the program.

<img width="735" alt="Screenshot 2025-04-26 at 5 32 19 PM" src="https://github.com/user-attachments/assets/62846b13-8567-429c-b93d-91b5625eb6e8" />
<img width="717" alt="Screenshot 2025-04-26 at 5 36 50 PM" src="https://github.com/user-attachments/assets/52ee8da9-f1f7-4bed-8961-eb610377f987" />

---

## Prerequisites

- macOS with Homebrew installed
- Basic knowledge of terminal commands
- Go installed on your macOS host (`brew install go`)

---

## Step 1: Install Lima

1. Install Lima using Homebrew:
   ```bash
   brew install lima
   ```

2. Start the default Lima instance:
   ```bash
   limactl start
   ```

3. Verify that the instance is running:
   ```bash
   limactl list
   ```

   You should see an output similar to:
   ```
   NAME       STATUS     SSH                ARCH       CPUS    MEMORY    DISK      DIR
   default    Running    127.0.0.1:60022    aarch64    4       4GiB      100GiB    ~/.lima/default
   ```

---

## Step 2: Generate an SSH Key Pair

1. Generate an SSH key pair on your macOS host:
   ```bash
   ssh-keygen -t ecdsa -f ~/.lima/default/id_ecdsa -q -N ""
   ```

2. Verify that the key pair was created:
   ```bash
   ls ~/.lima/default/id_ecdsa*
   ```

---

## Step 3: Copy the Public Key to the Lima VM

1. Log into the Lima VM:
   ```bash
   limactl shell default
   ```

2. Create the `.ssh` directory and set the correct permissions:
   ```bash
   mkdir -p ~/.ssh
   chmod 700 ~/.ssh
   ```

3. Append the public key to the `authorized_keys` file:
   ```bash
   echo "<contents of ~/.lima/default/id_ecdsa.pub>" >> ~/.ssh/authorized_keys
   chmod 600 ~/.ssh/authorized_keys
   ```

   Replace `<contents of ~/.lima/default/id_ecdsa.pub>` with the actual contents of the public key file:
   ```bash
   cat ~/.lima/default/id_ecdsa.pub
   ```

4. Exit the Lima shell:
   ```bash
   exit
   ```

---

## Step 4: Configure SSH in the Lima VM

1. Log back into the Lima VM:
   ```bash
   limactl shell default
   ```

2. Edit the SSH server configuration file:
   ```bash
   sudo nano /etc/ssh/sshd_config
   ```

3. Ensure the following lines are present and uncommented:
   ```
   PubkeyAuthentication yes
   AuthorizedKeysFile .ssh/authorized_keys
   ```

4. Restart the SSH service:
   ```bash
   sudo systemctl restart ssh
   ```

5. Exit the Lima shell:
   ```bash
   exit
   ```

---

## Step 5: SSH into the VM with the Correct User

1. SSH into the VM using the correct username (`radhakrishna`) and private key:
   ```bash
   ssh -i ~/.lima/default/id_ecdsa radhakrishna@127.0.0.1 -p 60022
   ```

---

## Step 6: Transfer the Go Program to the VM

1. Use `scp` to copy the `main.go` file to the VM:
   ```bash
   scp -P 60022 -i ~/.lima/default/id_ecdsa /Users/radhakrishna/GolandProjects/Software-Engineering-In-Depth-1/Networking/EBPF/main.go radhakrishna@127.0.0.1:/home/radhakrishna.linux/
   ```

2. Verify that the file was successfully transferred:
   ```bash
   ssh -i ~/.lima/default/id_ecdsa radhakrishna@127.0.0.1 -p 60022
   ls -l /home/radhakrishna.linux/main.go
   ```

---

## Step 7: Run the Go Program in the VM

1. Log into the VM:
   ```bash
   ssh -i ~/.lima/default/id_ecdsa radhakrishna@127.0.0.1 -p 60022
   ```

2. Navigate to the directory containing the `main.go` file:
   ```bash
   cd /home/radhakrishna.linux/
   ```

3. Run the Go program with `sudo`:
   ```bash
   sudo go run main.go
   ```

---

## Example Output

When the program runs successfully, you should see output similar to the following:
```
Running eBPF program to monitor active connections...
PID: 1234, Comm: curl, Port: 443
PID: 5678, Comm: ssh, Port: 22
...
```

---

## Troubleshooting

1. **Permission Denied (publickey)**:
   - Ensure the public key in the VM's `~/.ssh/authorized_keys` matches the private key on your macOS host.

2. **bpftrace Requires Root**:
   - Always run the program with `sudo` to ensure `bpftrace` has the necessary permissions.

3. **Missing Dependencies**:
   - Install `bpftrace` and `golang` in the VM:
     ```bash
     sudo apt update
     sudo apt install -y bpftrace golang
     ```

4. **Check SSH Logs**:
   - If SSH issues persist, check the server logs in the VM:
     ```bash
     sudo journalctl -u ssh
     ```

---

## Conclusion

You have successfully set up a Lima VM, configured SSH, transferred the Go program, and executed an eBPF program to monitor active network connections. Feel free to extend the program to capture additional details or monitor specific events.
