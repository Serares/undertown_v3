import { Response, Request } from 'express';

export const renderHomepage = (req: Request, res: Response) => {
    return res.render('pages/home/home', {
        pageTitle: 'Test Views',
        path: '/',
    });
};
