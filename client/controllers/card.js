Template.card.created = function() {
    // TODO: Ugly hack, make it beautiful and safe
    var cardType = this.view.parentView.name.split(".")[1];
    if(Template[cardType + "Details"] === undefined) {
        console.error("Card details template '" + cardType + "Details', doesn't exist");
    } else {
        TemplateVar.set(this, "cardType", cardType);
    }
};

Template.card.helpers({
    error: "",
    showDetails: function() {
        return this.singleCard ? true : TemplateVar.get("showDetails");
    },
    detailsTemplate: function() {
        return Template[TemplateVar.get("cardType") + "Details"];
    },
    cardType: function() {
        return TemplateVar.get("cardType");
    }
});

Template.card.events({
    "click #showDetails": function(e, template) {
        TemplateVar.set(template, "showDetails", !TemplateVar.get("showDetails"));
    }
});
