import { M3eSnackbar } from '@m3e/web/snackbar';

export function notify(message: string): void {
  M3eSnackbar.open(message);
}
