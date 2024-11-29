# Helheim
Second attempt at a C2 written in GO, MINIMAL CHATGPT HELP (ALMOST NONE)


Okay TEMP notes:

Server certs:
- Edit server.cnf and then compile with this:
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout server.key \
    -out server.crt \
    -config server.cnf
```
- Once done move server.crt to the clients. I'll figure out a better way to do this. I want to spend my energy on actually writing out the functionailty first. 