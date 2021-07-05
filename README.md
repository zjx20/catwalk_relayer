catwalk_relayer
===============

# Deploy To IBM Cloud Foundry

HomePage: https://cloud.ibm.com/cloudfoundry/overview

Getting Start: https://cloud.ibm.com/docs/cloud-foundry-public?topic=cloud-foundry-public-getting-started-go

Docs: https://docs.cloudfoundry.org/devguide/index.html

Note: not accessible directly

Deploy:
```bash
cd catwalk_relayer

ibmcloud cf login

# replace <app-name> to the real one
ibmcloud cf create-app-manifest <app-name>
ibmcloud cf push <app-name> -c "catwalk_relayer -upstream example.com:80 -ws true"
```

# Deploy To Heroku

HomePage: https://dashboard.heroku.com/

Getting Start: https://devcenter.heroku.com/articles/getting-started-with-go

Deploy:
```bash
heroku login

heroku git:remote -a <appName>
heroku config:set UPSTREAM=youserver.com:1234 WS=true WS_PATH=/chat -a <appName>

git push heroku master
```
