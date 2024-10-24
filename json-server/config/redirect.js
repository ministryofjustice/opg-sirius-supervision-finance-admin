/**
 * Redirect middleware to mock auth
 */
module.exports = (req, res, next) => {
    if (req.method === "GET" && req.query.redirect) {
        console.log("redirecting to " + req.query.redirect);
        res.cookie("sirius", "session");
        res.redirect(req.query.redirect);
    }
    next();
};
