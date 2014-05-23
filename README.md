Go Say Hello!
==============
Openshift Live: [http://gosayhello-stujo.rhcloud.com/](http://gosayhello-stujo.rhcloud.com/)
Heroku Live: [http://gosayhello.herokuapp.com/](http://gosayhello.herokuapp.com/)


#Installation
==============

##Via Home Brew

    brew install go --cross-compile-common

##Set up GOPATH

I chose to create customized Terminal Window

(GOPATH cannot start with ~)

Shell runs ```. ~/goshell.sh``` when opened:

    #!/bin/sh
    export GOPATH="/usr/local/var/go"
    export PATH=$PATH:$GOPATH/bin
    cd $GOPATH/src

#Basic Modules
It looks like go get uses mercurial, so I needed this

    $ brew install mercurial

Heroku Deployment needs this?

    go get github.com/tools/godep

##Check the Setup

    $ which go
    /usr/local/bin/go
    $ go version
    go version go1.2.2 darwin/amd64


##Trying to run the app?

    $ go run web.go
    listening on :...panic: listen tcp: unknown port tcp/

:(

##Ahhh Environment Variables

The code is setup to run on openshift where PORT and HOST are defined in the environment

    export PORT=3337
    export HOST=localhost

And voila!

    $ go run web.go 
    listening on localhost:3337...

![Hello World][hello-world-screenshot.png]

#Local File Paths

After a bit of experimenting the layout that worked for me was to create the project folder
under $GOLANG/src

Specifically, matching the entire path specified in the .godir

```github.com/stujo/gosayhello```

to give a full project path of

```/usr/local/var/go/src/github.com/stujo/gosayhello```

or

```$GOPATH/src/github.com/stujo/gosayhello```


#Heroku

Deployment Failed:

    Warning: Permanently added 'heroku.com,50.19.85.156' (RSA) to the list of known hosts.
    Initializing repository, done.
    Counting objects: 33, done.
    Delta compression using up to 32 threads.
    Compressing objects: 100% (27/27), done.
    Writing objects: 100% (33/33), 1.69 MiB | 0 bytes/s, done.
    Total 33 (delta 11), reused 0 (delta 0)

     !     Push rejected, no Cedar-supported app detected
    To git@heroku.com:gosayhello.git
     ! [remote rejected] 0470adb8cb56b62754fca2d06e8551e66cfae812 -> master (pre-receive hook declined)
    error: failed to push some refs to 'git@heroku.com:gosayhello.git'


##Connect the app

    heroku git:remote -a gosayhello

##Set the buildpack

    heroku config:set BUILDPACK_URL=https://github.com/kr/heroku-buildpack-go.git

##Not specifying the HOSTNAME on the port binding

I found that I had to add code to NOT specify a hostname in the port binding on HEROKU,
this is probably because I originally had HOST specified in the config. I added ON_HEROKU
to detect heroku, so if you run into these issues, and are not trying to deploy to
multiple platforms just change the code to not specify a HOST in any environment

    $ heroku config
    === gosayhello Config Vars
    BUILDPACK_URL: https://github.com/kr/heroku-buildpack-go.git
    ON_HEROKU:     1

##Deploy again

    2014-05-23T05:46:46.686083+00:00 app[web.1]: Starting Go Say Hello
    2014-05-23T05:46:46.686086+00:00 app[web.1]: HOST ()
    2014-05-23T05:46:46.686087+00:00 app[web.1]: PORT (5414)
    2014-05-23T05:46:45.615447+00:00 heroku[web.1]: Starting process with command `gosayhello`
    2014-05-23T05:46:47.379438+00:00 heroku[web.1]: State changed from starting to up

yeay!


#Openshift Deployment

(I renamed my git remote for openshift from origin to openshift)


    $ git push openshift master
    Counting objects: 5, done.
    Delta compression using up to 4 threads.
    Compressing objects: 100% (3/3), done.
    Writing objects: 100% (3/3), 1011 bytes | 0 bytes/s, done.
    Total 3 (delta 2), reused 0 (delta 0)
    remote: Building git ref 'master', commit aa46cc8
    remote:
    remote: -----> Using Go 1.2.2
    remote: -----> Running: go get -tags openshift ./...
    remote: Preparing build for deployment
    remote: Deployment id is edade3bb
    remote: Activating deployment
    remote:
    remote: -------------------------
    remote: Git Post-Receive Result: success
    remote: Activation status: success
    remote: Deployment completed with status: success
    To ssh://537e8c40e0b8cd3eb80002e0@gosayhello-stujo.rhcloud.com/~/git/gosayhello.git/
       091c90c..aa46cc8  master -> master

yeay!




------------------------


Original:

OpenShift Go Cartridge
======================

Runs [Go](http://golang.org) on [OpenShift](https://openshift.redhat.com/app/login) using downloadable cartridge support. 

Once the app is created, you'll have a ".godir" file in the root of your repo. The single line is to tell the cartridge what the package of your Go code is.  A typical .godir file might contain:

    github.com/stujo/gosayhello

which would tell OpenShift to place all of the files in the root of the Git repository inside of the <code>github.com/stujo/gosayhello</code> package prior to compilation.

When you push code to the repo, the cart will compile your package into <code>$OPENSHIFT_REPO_DIR/bin/</code>, with the last segment of the .godir being the name of the executable.  For the above .godir, your executable will be:

    $OPENSHIFT_REPO_DIR/bin/gosayhello

If you want to serve web requests (vs. running in the background), you'll need to listen on the ip address and port that OpenShift allocates - those are available as HOST and PORT in the environment.

This default "web.go" file is a simple "hello, world" web service. 

Any log output will be generated to <code>$OPENSHIFT_GO_LOG_DIR</code> on your OpenShift gear


Build
-----

When you push code to your repo, a Git postreceive hook runs and invokes a compile script.  This attempts to download the Go compiler environment for you into $OPENSHIFT_GO_DIR/cache.  Once the environment is setup, the cart runs

    go get -tags openshift ./...

on a working copy of your source. 
The main file that you run will have access to two environment variables, $HOST and $PORT, which contain the internal address you must listen on to receive HTTP requests to your application.

