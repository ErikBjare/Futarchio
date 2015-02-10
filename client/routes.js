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
Router.route('/about', function() {
    this.render("about");
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
            data = Statements.findOne({"_id": id});
            if(data === undefined) {
                return data;
            }
            data.singleCard = true;
            data.cardType = "statement";
            return data;
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
            var user = Meteor.users.findOne({"username": username});
            return user;
        }
    });
});
