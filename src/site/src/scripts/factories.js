app.factory('msgStack', function() {
    var msgs = [];
    return function(obj) {
        if(msg === "") return msgs;
        msgs.push({type: "info", msg: msg});
        console.log(obj);
    };
});

app.factory('Poll', function($resource, user) {
    return $resource("/api/0/polls", {},
        {"save": {method: "POST", isArray: false, headers: {"Authorization": user.authkey()}}});
});

app.factory('Vote', function($log, $resource, user) {
    return $resource('/api/0/polls/:pollid/vote', {},
        {"save": {method: "POST", headers: {"Authorization": user.authkey()}}});
});

app.factory('Statement', function($resource, user) {
    return $resource("/api/0/statements", {},
        {"save": {method: "POST", isArray: false, headers: {"Authorization": user.authkey()}}});
});

app.factory('Prediction', function($log, $resource, user) {
    return $resource('/api/0/statements/:key/predict', {},
        {"save": {method: "POST", headers: {"Authorization": user.authkey()}}});
});

app.factory('User', function($log, $resource) {
    return $resource('/api/0/users', {}, {});
});

app.factory('gravatar', function($log, $resource) {
    var gravatar = {};
    gravatar.hash = function(email) {
        if(!_.contains(email, "@")) {
            console.error("Error: got empty email");
        }
        var hash = CryptoJS.MD5(email).toString();
        return hash;
    };

    return gravatar;
});

app.factory('user', function($q, $log, $http, $route, $cookieStore, $location, $window, gravatar) {
    var user = {};

    user.notifications = [];
    /*
    user.notifications = [
        {title: "Example notification",
         description: "Example description, lalalalala.\nasdasdad",
         url: "/notifications/1"},
        {title: "Notif 2",
         description: "Another one",
         url: "/notifications/2"}
    ];
    */

    // Returns true if logged in, else false
    user.is_logged_in = function() {
        val = $cookieStore.get("me") && user.authkey();
        return val ? true : false;
    };

    // Getter and setter for authkey
    user.authkey = function(authkey) {
        if(authkey !== undefined) {
            $cookieStore.put("auth", authkey);
            console.log("authkey cookie put");
        } else {
            authkey = $cookieStore.get("auth");
        }
        return authkey;
    };

    // Attempts to log in with username and password
    user.login = function(username, password) {
        var deferred = $q.defer();

        $http.post('/api/0/auth', {username: username, password: password})
            .success(function(data, status, headers, config) {
                console.log(data);
                user.authkey(data.key);
                $http.get('/api/0/users/me', {"headers": {"Authorization": user.authkey()}})
                .success(function(data) {
                    $cookieStore.put("me", data);
                    deferred.resolve(data);
                }).error(function(data) {
                    deferred.reject("error while fetching profile");
                });
            }).error(function(data, status, headers, config) {
                deferred.reject(data.error);
            });

        return deferred.promise;
    };

    user.logout = function() {
        $cookieStore.remove("auth");
        $cookieStore.remove("me");
        $location.path("/");

        // TODO: Use $window.location.reload() instead?
        $route.reload();
    };

    user.username = function() {
        return $cookieStore.get("me").username;
    };

    user.email = function() {
        return $cookieStore.get("me").email;
    };

    user.gravatar_hash = function() {
        return gravatar.hash(user.email());
    };

    user.votereceipts = undefined;
    user.get_votereceipts = function() {
        $http.get('/api/0/polls/myvotereceipts', {'headers': {'Authorization': user.authkey()}})
            .success(function(data) {
                user.votereceipts = data;
            }).error(function(data) {
                $log.error(data);
            });
    };
    if(user.is_logged_in()) {
        user.get_votereceipts();
    }


    user.has_voted_on = function(pollid) {
        if(user.votereceipts === undefined) {
            console.log("Votereceipts not received");
            return false;
        }
        for(var i = 0; i < user.votereceipts.length; i++) {
            if(user.votereceipts[i].pollid == pollid) {
                return true;
            }
        }
        return false;
    };

    return user;
});
