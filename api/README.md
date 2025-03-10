# API


Example of task running
```bash
curl -X POST http://localhost:8080/api/task \
     -H "Content-Type: application/json" \
     -d '{
           "type": "ansible",
           "params": {
             "inventory": "my_inventory.ini",
             "playbook": "site.yml"
           }
         }'

```