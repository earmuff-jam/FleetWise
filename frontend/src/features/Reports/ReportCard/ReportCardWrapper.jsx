import { Card, CardContent, Chip, Stack, Typography } from '@mui/material';

export default function ReportCardWrapper({ title, chipLabel, value, footerText, footerSuffix, dataTour, children }) {
  return (
    <Card
      data-tour={dataTour}
      sx={{
        display: 'flex',
        flexDirection: 'column',
        borderRadius: '0.2rem',
        width: { xs: '100%', sm: '100%' },
      }}
    >
      <CardContent>
        <Stack direction="row">
          <Stack flexGrow={1} sx={{ minWidth: '12rem', minHeight: '6rem' }}>
            <Stack direction="row" justifyContent="space-between">
              <Typography variant="h6" component="h3">
                {title}
              </Typography>
              {chipLabel && <Chip label={chipLabel} variant="outlined" sx={{ borderRadius: '0.3rem' }} />}
            </Stack>
            <Stack direction="row" alignItems="center" justifyContent="space-between">
              {value && (
                <Typography variant="h1" flexWrap={1} color="text.secondary">
                  {value}
                </Typography>
              )}
            </Stack>
            {children}
          </Stack>
        </Stack>
      </CardContent>

      {footerText && (
        <Stack direction="row" spacing={0.3} paddingLeft={1}>
          <Typography variant="caption" color="text.secondary">
            {footerText}
          </Typography>
          <Typography variant="caption" component={'span'} color="success.main">
            {footerSuffix}
          </Typography>
        </Stack>
      )}
    </Card>
  );
}
