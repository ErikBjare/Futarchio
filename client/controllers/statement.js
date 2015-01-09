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
        var pred = Prediction({
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
        statements = Statements.find({}, {sort: {"createdAt": -1}}).fetch();
        statements = _.map(statements, function(obj) {
            return _.extend(obj, {"cardType": "statement"});
        });
        return statements;
    }
});

Template.statements.events({
    "click button#addbtn": function(event) {
        Session.set("showNewStmt", !Session.get("showNewStmt"));
    }
});

Template.newstatement.events({
    "submit": function(event) {

        // Empty variable field handling
        if (!event.target.title.value){
            Session.set("formerror", "You cannot create a statement without a title");
            return false;
        }
        else if (!event.target.description.value){
            Session.set("formerror", "You cannot create a statement without a description");
            return false;
        }

        var stmt = Statement({
            title: event.target.title.value,
             description: event.target.description.value
        });
        Statements.insert(stmt);

        event.target.title.value = "";
        event.target.description.value = "";

        return false;
    }
});

