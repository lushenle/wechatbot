#!/bin/bash
rsync -avzP --exclude=".idea" --exclude=".git" --delete ~/workspace/wechatbot/ $1:~/workspace/wechatbot/
