app.factory('msgStack', function() {
    var msgs = [];
    return function(obj) {
        if(msg === "") return msgs;
        msgs.push({type: "info", msg: msg});
        console.log(obj);
    };
});

app.factory('UserKeyVal', function($log, $resource) {
    // DEPRECATED
    // TODO: Remove refs
    var User = $resource('/api/0/users/:key/:val', {});
    return function(key, val) {
        var user = User.query(
           {"key": key,
            "val": val},
            function(u) {
                $log.info("Successfully fetched user");
            }, function(u) {
                $log.error("Error");
            });
        return user;
    };
});

app.factory('User', function($log, $resource) {
    return $resource('/api/0/users', {}, {});
});

app.factory('Vote', function($log, $resource, user) {
    return $resource('/api/0/polls/:pollid/vote?api_key='+user.getAuthkey());
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
    user.is_logged_in = function() {
        val = $cookieStore.get("me") && $cookieStore.get("auth");
        return val;
    };

    user.getAuthkey = function() {
        authkey = $cookieStore.get("auth");
        return authkey;
    };

    user.login = function(username, password) {
        var deferred = $q.defer();

        $http.post('/api/0/auth', {username: username, password: password})
        .success(function(data, status, headers, config) {
            console.log(data);
            authkey = data.key;
            $cookieStore.put("auth", authkey);
            console.log("Cookie put");
            $http.get('/api/0/users/me', {"headers": {"Authorization": authkey}})
            .success(function(data) {
                $cookieStore.put("me", {"username": data.username, "email": data.email});
                deferred.resolve(data);
            }).error(function(data) {
                deferred.reject("error while fetching profile");
            });
        }).error(function(data, status, headers, config) {
            deferred.reject(data.error)
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

    return user;
});
