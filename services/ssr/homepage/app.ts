import express from 'express';
import homepageRouter from './routes/home';

const app = express();
const port = process.env.PORT || 8080;

app.set('views', __dirname + '/views');
app.set('view engine', 'ejs');

app.use('/', homepageRouter);

app.listen(port, () => {
    // log server started
    console.log(`server started at localhost:${port}`);
});
