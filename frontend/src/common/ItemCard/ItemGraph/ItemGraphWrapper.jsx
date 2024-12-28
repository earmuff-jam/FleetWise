import ItemGraph from '@common/ItemCard/ItemGraph/ItemGraph';
import { Paper } from '@mui/material';

export default function ItemGraphWrapper({ associatedAssets = [], graphDataTour }) {
  return (
    <Paper elevation={1} sx={{ padding: '1rem' }} data-tour={graphDataTour}>
      <ItemGraph associatedAssets={associatedAssets} />
    </Paper>
  );
}
