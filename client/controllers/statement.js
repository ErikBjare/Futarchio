var defaultOrder = "date";

Template.statements.created = function() {
    TemplateVar.set(this, "order", defaultOrder);
};

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
        // TODO: Sort by points & activity by publishing and subscribing to aggregations server-side
        // https://stackoverflow.com/questions/18520567/average-aggregation-queries-in-meteor
        // Votes.aggregate([{$match: {type: "UpDown"}}, {$group: {_id: "$post", score: {$sum: "$value"}}}]);

        var orderToSort = {
            "date": {createdAt: -1},
            "points": {points: -1},
            "activity": {createdAt: -1}
        };
        statements = Statements.find({}, orderToSort).fetch();
        statements = _.map(statements, function(obj) {
            return _.extend(obj, {"cardType": "statement"});
        });
        return statements;
    },
    order: function() {
        return TemplateVar.get("order");
    }
});

Template.statements.events({
    "click button#newbtn": function(e, template) {
        TemplateVar.set(template, "showNewStmt", !TemplateVar.get(template, "showNewStmt"));
    }
});

Template.newstatement.events({
    "submit": function(event, template) {

        // Empty variable field handling
        if (!event.target.title.value){
            TemplateVar.set(template, "formerror", "You cannot create a statement without a title");
            return false;
        }
        else if (!event.target.description.value){
            TemplateVar.set(template, "formerror", "You cannot create a statement without a description");
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

