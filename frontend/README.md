# Literary Opinions Graph - Frontend

A Next.js frontend application for visualizing relationships between writers and literary works based on documented opinions.

## Tech Stack

- **Framework**: Next.js 16+ with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Testing**: React Testing Library (setup ready)

## Project Structure

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── layout.tsx          # Root layout with providers
│   │   └── page.tsx            # Home page
│   ├── components/             # React components
│   │   ├── writers/            # Writer-related components
│   │   ├── works/              # Work-related components
│   │   ├── opinions/           # Opinion-related components
│   │   └── common/             # Shared components (Button, Input, etc.)
│   ├── services/              # API service layer
│   │   ├── writerService.ts
│   │   ├── workService.ts
│   │   └── opinionService.ts
│   ├── stores/                 # Zustand stores
│   │   ├── writerStore.ts
│   │   ├── workStore.ts
│   │   └── opinionStore.ts
│   ├── types/                  # TypeScript type definitions
│   │   ├── writer.ts
│   │   ├── work.ts
│   │   └── opinion.ts
│   └── hooks/                  # Custom React hooks
├── public/                     # Static assets
└── .env.example                # Example environment variables
```

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn/pnpm
- Backend API running (see backend README)

### Installation

1. Install dependencies:
```bash
npm install
# or
yarn install
# or
pnpm install
```

2. Set up environment variables:
```bash
cp .env.example .env.local
```

3. Update `.env.local` with your backend API URL:
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Development

Start the development server:
```bash
npm run dev
# or
yarn dev
# or
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Create production build
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm run lint:fix` - Run ESLint and auto-fix issues
- `npm run format` - Format code with Prettier
- `npm run format:check` - Check code formatting
- `npm run type-check` - Run TypeScript type checking
- `npm run check` - Run all checks (type-check, lint, format-check)

## Environment Variables

| Variable              | Description          | Default                        |
| --------------------- | -------------------- | ------------------------------ |
| `NEXT_PUBLIC_API_URL` | Backend API base URL | `http://localhost:8080/api/v1` |

## API Integration

The frontend communicates with the Go backend API at `/api/v1/` endpoints:

- **Writers**: `/api/v1/writers` - CRUD operations
- **Works**: `/api/v1/works` - CRUD operations + GetByAuthor
- **Opinions**: `/api/v1/opinions` - CRUD operations + GetByWriter, GetByWork, GetByWriterAndWork

All API calls are handled through service classes in `src/services/` and state is managed via Zustand stores in `src/stores/`.

## Code Standards

This project follows strict coding standards defined in `.cursor/rules/frontend.mdc`:

- Functional components with TypeScript interfaces
- Named exports (not default exports)
- Explicit return types for exported functions
- Proper error and loading state handling
- Semantic HTML and accessibility
- Small, focused components with single responsibility

## Linting and Formatting

The project uses ESLint and Prettier for code quality and consistency:

### ESLint Configuration

- **TypeScript**: Strict type checking, no `any` types, explicit return types
- **React**: Hooks rules, accessibility checks
- **Code Quality**: No console/debugger, prefer const, template literals
- **Accessibility**: JSX a11y rules enforced

### Prettier Configuration

- **Print Width**: 100 characters
- **Tab Width**: 2 spaces
- **Semicolons**: Required
- **Trailing Commas**: ES5 style
- **Single Quotes**: Disabled (double quotes)

### Running Linters

```bash
# Check for linting issues
npm run lint

# Auto-fix linting issues
npm run lint:fix

# Check code formatting
npm run format:check

# Format all code
npm run format

# Run all checks (type-check, lint, format-check)
npm run check
```

All code should pass linting and formatting checks before committing.

## Development Workflow

1. Create components following the component template pattern
2. Use Zustand stores for state management
3. Use service classes for API calls
4. Handle loading and error states in all async operations
5. Write tests focusing on user interactions

## Testing

Testing setup is ready for:
- Component tests with React Testing Library
- Service tests with mocked fetch
- Store tests with Zustand testing utilities

## Learn More

- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Zustand Documentation](https://zustand-demo.pmnd.rs/)
- [TypeScript Documentation](https://www.typescriptlang.org/docs/)
