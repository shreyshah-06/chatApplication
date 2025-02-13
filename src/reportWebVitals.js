export function reportWebVitals(onPerfEntry) {
  if (onPerfEntry && onPerfEntry instanceof Function) {
    import('web-vitals').then(({ getCLS, getFID, getLCP }) => {
      // Logs Cumulative Layout Shift (CLS)
      getCLS(onPerfEntry);
      // Logs First Input Delay (FID)
      getFID(onPerfEntry);
      // Logs Largest Contentful Paint (LCP)
      getLCP(onPerfEntry);
    });
  }
}
