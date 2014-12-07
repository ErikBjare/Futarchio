Template.polls.helpers({
    polls: function() {
        console.log("Fetching polls");
        polls = Polls.find({}, {sort: {createdAt: -1}});
        return polls;
    }
});

Template.poll.helpers({
    error: ""
});

Template.statement.events({
    "click button#predict": function() {
        console.log(this);
        Predictions.insert({
            credence: event.target.credence.value,
            createdBy: Meteor.userId()
        });
        console.log("Inserted prediction");
    }
});

Template.statements.helpers({
    statements: function() {
        console.log("Fetching statements");
        stmts = Statements.find({}, {order: {"createdAt": -1}});
        return stmts;
    },
    showAdd: function() {
        return Session.get("addStmt-visible");
    }
});

Template.statements.events({
    "click button#addbtn": function(event) {
        Session.set("addStmt-visible", !Session.get("addStmt-visible"));
    }
});

Template.newpoll.events({
    "submit": function(event) {
        Polls.insert({
            title: event.target.title.value,
            description: event.target.description.value,
            createdBy: Meteor.userId(),
            createdAt: new Date()
        });

        event.target.title.value = "";
        event.target.description.value = "";

        return false;
    }
});

Template.newstatement.events({
    "submit": function(event) {
        Statements.insert({
            title: event.target.title.value,
            description: event.target.description.value,
            createdBy: Meteor.userId(),
            createdAt: new Date()
        });

        event.target.title.value = "";
        event.target.description.value = "";

        return false;
    }
});
