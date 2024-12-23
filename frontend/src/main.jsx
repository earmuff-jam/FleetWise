import { store } from '@src/Store';
import { Provider } from 'react-redux';
import { lightTheme } from '@utils/Theme';

import * as ReactDOM from 'react-dom/client';
import { SnackbarProvider } from 'notistack';

import { CssBaseline, ThemeProvider } from '@mui/material';
import ApplicationValidator from '@src/ApplicationValidator';

ReactDOM.createRoot(document.getElementById('root')).render(
  <ThemeProvider theme={lightTheme}>
    <CssBaseline />
    <Provider store={store}>
      <SnackbarProvider
        dense
        preventDuplicate
        maxSnack={3}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
        autoHideDuration={3000}
      >
        <ApplicationValidator />
      </SnackbarProvider>
    </Provider>
  </ThemeProvider>
);
