const EXIT_STYLE = `
.base:not([open]) {
  display: none !important;
  transition:
    transform 200ms cubic-bezier(0.3, 0, 0.8, 0.15),
    overlay 200ms allow-discrete,
    visibility 200ms allow-discrete,
    display 200ms allow-discrete !important;
}
@media (prefers-reduced-motion: reduce) {
  .base:not([open]) {
    transition: none !important;
  }
}
`;

function paintDialogExit(dialog: Element, attempt = 0): void {
  const root = dialog.shadowRoot;
  if (!root) {
    if (attempt < 5 && dialog.isConnected) {
      requestAnimationFrame(() => paintDialogExit(dialog, attempt + 1));
    }
    return;
  }
  let style = root.getElementById('yukumo-dialog-motion');
  if (!style) {
    style = document.createElement('style');
    style.id = 'yukumo-dialog-motion';
    root.append(style);
  }
  if (style.textContent !== EXIT_STYLE) style.textContent = EXIT_STYLE;
}

function scan(node: ParentNode): void {
  if (node instanceof Element && node.localName === 'm3e-dialog') {
    paintDialogExit(node);
  }
  if (node instanceof Element || node instanceof Document || node instanceof DocumentFragment) {
    node.querySelectorAll('m3e-dialog').forEach(paintDialogExit);
  }
}

export function installDialogMotion(): void {
  scan(document);
  const observer = new MutationObserver((records) => {
    for (const record of records) {
      record.addedNodes.forEach((node) => {
        if (node instanceof Element) scan(node);
      });
    }
  });
  observer.observe(document.documentElement, { childList: true, subtree: true });
}

export function waitForDialogExit(dialog: HTMLElement): Promise<void> {
  const base = dialog.shadowRoot?.querySelector('.base');
  if (!(base instanceof HTMLElement)) return Promise.resolve();
  return new Promise((resolve) => {
    let settled = false;
    const done = (): void => {
      if (settled) return;
      settled = true;
      base.removeEventListener('transitionend', onEnd);
      resolve();
    };
    const onEnd = (event: TransitionEvent): void => {
      if (event.target !== base) return;
      if (event.propertyName !== 'transform' && event.propertyName !== 'display') {
        return;
      }
      done();
    };
    base.addEventListener('transitionend', onEnd);
    window.setTimeout(done, 280);
  });
}
