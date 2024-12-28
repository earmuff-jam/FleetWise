import { Box, Button, Link, Stack, Typography } from '@mui/material';

const RowHeader = ({
  title,
  titleVariant = 'h4',
  caption,
  showRedirectLink,
  redirectTo,
  primaryButtonTextLabel,
  primaryButtonDisabled,
  primaryStartIcon,
  secondaryButtonTextLabel,
  secondaryStartIcon,
  secondaryButtonDisabled,
  handleClickPrimaryButton,
  handleClickSecondaryButton,
  primaryBtnDataTour,
  secondaryBtnDataTour,
  children,
}) => {
  return (
    <Box display="flex" justifyContent="space-between" alignItems="center">
      <Stack>
        <Typography variant={titleVariant} component="h2" color="text.secondary">
          {title}
        </Typography>
        {caption ? (
          showRedirectLink ? (
            <Link href={redirectTo}>
              <Typography variant="caption">{caption}</Typography>
            </Link>
          ) : (
            <Typography variant="caption">{caption}</Typography>
          )
        ) : null}
      </Stack>
      <Stack direction="row" spacing={2} useFlexGap>
        {primaryButtonTextLabel ? (
          <Button
            size="small"
            color="primary"
            variant="outlined"
            data-tour={primaryBtnDataTour}
            disabled={primaryButtonDisabled}
            onClick={handleClickPrimaryButton}
            startIcon={primaryStartIcon}
          >
            {primaryButtonTextLabel}
          </Button>
        ) : null}
        {secondaryButtonTextLabel ? (
          <Button
            size="small"
            color="primary"
            variant="outlined"
            data-tour={secondaryBtnDataTour}
            disabled={secondaryButtonDisabled}
            onClick={handleClickSecondaryButton}
            startIcon={secondaryStartIcon}
          >
            {secondaryButtonTextLabel}
          </Button>
        ) : null}
        {children}
      </Stack>
    </Box>
  );
};

export default RowHeader;
