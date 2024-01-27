import React from 'react';
import { Button, TextField } from '@mui/material';

export default function SubmitProperty() {
    return (
      <form>
        <TextField id="title" label="Title" variant="outlined" fullWidth margin="normal" />
        <input type="file" />
        <Button variant="contained" color="primary" type="submit">
          Submit
        </Button>
      </form>
    );
}
