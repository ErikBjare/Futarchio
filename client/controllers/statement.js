Template.statementDetails.helpers({
    predictions: function() {
        return Predictions.find({"statement": this._id}, {sort: {"createdAt": -1}});
    },
    predictionsCount: function() {
        return Predictions.find({"statement": this._id}).count();
    }
});

Template.statementDetails.events({
    "submit": function(event, template) {
        console.log(this);
        console.log(this.post);
        var pred = new Prediction({
            "credence": event.target.credence.value,
            "statement": this._id
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

Template.newstatement.events({
    "submit": function(event) {
        var stmt = new Statement({
            title: event.target.title.value,
             description: event.target.description.value
        });
        Statements.insert(stmt);

        event.target.title.value = "";
        event.target.description.value = "";

        return false;
    }
});

