## what?
Gofocus disables access too websites for a given amount of time.

Add websites you want to be blocked in `hosts_config.json`:
```json
{
  "ip_address": "127.0.0.1",
  "host_names": [
    {
      "host_name": "www.youtube.com"
    },
    {
      "host_name": "https://news.ycombinator.com/"
    }
  ]
}
```
## todo
- create installer
- workaround for sudo
- tbd
<img width="584" alt="Screenshot 2021-12-29 at 16 30 44" src="https://user-images.githubusercontent.com/44348300/147678024-98132b81-ed03-46a9-8034-1a410c4e0560.png">

Now go focus.
