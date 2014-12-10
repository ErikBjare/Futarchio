Accounts.ui.config({
    requestPermissions: {
      facebook: ['user_likes'],
      github: ['user', 'repo']
    },
    requestOfflineToken: {
      google: true
    },
    passwordSignupFields: 'USERNAME_AND_EMAIL'
});

Template.polls.helpers({
    polls: function() {
        return Polls.find({}, {sort: {createdAt: -1}});
    }
});

Template.polls.events({
    "click #newPollBtn": function() {
        Session.set("showNewPoll", !Session.get("showNewPoll"));
    }
});

Template.poll.created = function() {
    this.data.showDetails = new ReactiveVar();
    this.data.showDetails.set(false);
};

Template.poll.helpers({
    error: "",
    showDetails: function() {
        return this.showDetails.get();
    }
});

Template.poll.events({
    "click #showDetails": function(e, template) {
        template.data.showDetails.set(!(template.data.showDetails.get()));
    }
});

Template.statement.helpers({
    predictions: function() {
        return Predictions.find({"statement": this._id}, {sort: {"createdAt": -1}});
    },
    predictionsCount: function() {
        return Predictions.find({"statement": this._id}).count();
    }
});

Template.statement.events({
    "submit": function(event, template) {
        pred = new Prediction({
            "credence": event.target.credence.value,
            "statement": template.data._id
        });
        Predictions.insert(pred);

        event.target.credence.value = "";

        return false;
    }
});

Template.statements.helpers({
    statements: function() {
        return Statements.find({}, {sort: {"createdAt": -1}});
    }
});

Template.statements.events({
    "click button#addbtn": function(event) {
        Session.set("showNewStmt", !Session.get("showNewStmt"));
    }
});

Template.newpoll.events({
    "submit": function(event) {
        poll = new Poll({
            title: event.target.title.value,
            description: event.target.description.value
        });
        Polls.insert(poll);

        event.target.title.value = "";
        event.target.description.value = "";

        return false;
    }
});

Template.newstatement.events({
    "submit": function(event) {
        stmt = new Statement({
            title: event.target.title.value,
            description: event.target.description.value
        });
        Statements.insert(stmt);

        event.target.title.value = "";
        event.target.description.value = "";

        return false;
    }
});

Template.votebuttons.helpers({
    score: function() {
        upVotes = Votes.find({"post": this._id, value: 1}).count();
        downVotes = Votes.find({"post": this._id, value: -1}).count();
        score = upVotes - downVotes;
        return score.toString();
    }
});

Template.votebuttons.events({
    "click #up": function(event, template) {
        vote = new Vote({
            type: "UpDown",
            value: 1,
            post: template.data._id
        });
        Votes.insert(vote);
    },
    "click #down": function(event, template) {
        vote = new Vote({
            type: "UpDown",
            value: -1,
            post: template.data._id
        });
        Votes.insert(vote);
    }
});

