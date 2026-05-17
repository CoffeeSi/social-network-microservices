import '@/index.css';
import { Main } from '@/main';
import { App } from '@/App';

export const metadata = {
  title: 'Social Network',
  description: 'Client for microservices and API gateway',
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <Main>
          <App>{children}</App>
        </Main>
      </body>
    </html>
  );
}
