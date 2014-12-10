Router.configure({
    layoutTemplate: 'layout',
    loadingTemplate: 'loading',
    notFoundtemplate: 'notFound',

    // Track pageviews
    trackPageView: true
});

// Home
Router.route('/', function() {
    this.render("home");
});

// Votes and Polls
Router.route('/vote', function() {
    this.render("polls");
});

Router.route('/poll/:_id', function () {
    var params = this.params;
    var id = params._id;
    this.render("singlePoll", {
        data: function() {
            return Polls.findOne({"_id": id});
        }
    });
});

// Predictions and Statements
Router.route('/predict', function() {
    this.render("statements");
});

Router.route('/statement/:_id', function () {
    var params = this.params;
    var id = params._id;
    this.render("singleStatement", {
        data: function() {
            return Statements.findOne({"_id": id});
        }
    });
});

Router.route('/prediction/:_id', function () {
    var params = this.params;
    var id = params._id;
    this.render("prediction", {
        data: function() {
            return Predictions.findOne({"_id": id});
        }
    });
});

// User profiles
Router.route('/profile/:username', function () {
    var params = this.params;
    var username = params.username;
    this.render("profile", {
        data: function() {
            user = Meteor.users.findOne({"username": username});
            console.log(user);
            return user;
        }
    });
});
