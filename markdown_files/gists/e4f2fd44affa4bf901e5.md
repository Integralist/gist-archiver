# If a shell command returns an exit status then you'll find that in a CI environment that the job will immediately fail. To resolve this you should capture the error and then pipe the result and then use PIPESTATUS instead

[View original Gist on GitHub](https://gist.github.com/Integralist/e4f2fd44affa4bf901e5)

## pipestatus.sh

```shell
docker run \
  --rm \
  --cpu-shares 1024 \
  -v $WORKSPACE:/app \
  bbcnews/rubocop-config --fail-level F 2>&1 | true
    
RUBOCOP_RESULT="${PIPESTATUS[0]}"

# Similar, but different use case where we need to use $? to get result of last command

docker run \
  --rm \
  --cpu-shares 1024 \
  -v $WORKSPACE:/app \
  ruby:2.1.2 \
  bash -c "cd /app && gem install bundler && bundle install --jobs 4 && rspec" 2>&1

TESTS_RESULT="$?"
```

