import type { M3eButtonElement } from '@m3e/web/button';
import type { M3eDialogElement } from '@m3e/web/dialog';
import { waitForDialogExit } from './dialogMotion';

export interface ConfirmDialogOptions {
  title: string;
  message: string;
  confirmText: string;
  cancelText: string;
  danger?: boolean;
}

export function confirmDialog(options: ConfirmDialogOptions): Promise<boolean> {
  const dialog = document.createElement('m3e-dialog') as M3eDialogElement;
  dialog.alert = true;
  dialog.style.setProperty('--m3e-dialog-max-width', 'min(28rem, 92vw)');

  const header = document.createElement('span');
  header.slot = 'header';
  header.textContent = options.title;

  const message = document.createElement('p');
  message.textContent = options.message;

  const actions = document.createElement('div');
  actions.slot = 'actions';
  actions.toggleAttribute('end', true);

  actions.append(
    actionButton(options.cancelText, 'cancel', 'text'),
    actionButton(options.confirmText, 'confirm', 'filled', options.danger),
  );

  dialog.append(header, message, actions);

  return new Promise((resolve) => {
    const finish = (): void => {
      dialog.removeEventListener('closed', finish);
      resolve(dialog.returnValue === 'confirm');
      void waitForDialogExit(dialog).then(() => dialog.remove());
    };
    dialog.addEventListener('closed', finish);
    (document.querySelector('m3e-theme') ?? document.body).append(dialog);
    void dialog.show();
  });
}

function actionButton(
  label: string,
  returnValue: string,
  variant: M3eButtonElement['variant'],
  danger = false,
): M3eButtonElement {
  const button = document.createElement('m3e-button') as M3eButtonElement;
  button.variant = variant;
  if (danger) {
    button.classList.add('tone-error');
  }
  const action = document.createElement('m3e-dialog-action');
  action.setAttribute('return-value', returnValue);
  action.textContent = label;
  button.append(action);
  return button;
}
