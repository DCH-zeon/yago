<?php

namespace App\Http\Requests;

use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class PermissionRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     */
    public function authorize(): bool
    {
        return true;
    }

    protected function prepareForValidation(): void
    {
        $this->merge([
            'name' => trim($this->name, " /\\\n\r\t\v\0"),
            'route' => trim($this->route, " /\\\n\r\t\v\0"),
        ]);
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, ValidationRule|array|string>
     */
    public function rules(): array
    {
        return [
            'name' => [
                'required',
                'string',
                'max:255',
                Rule::unique('permissions')->where(function ($query) {
                    return $query->where('name', $this->name)
                        ->where('route', $this->route);
                })->ignore($this->permission),
            ],
            'route' => [
                'required',
                'string',
                'max:255',
                Rule::unique('permissions')->where(function ($query) {
                    return $query->where('route', $this->route)
                        ->where('name', $this->name);
                })->ignore($this->permission)
            ],
            'description' => 'nullable|string|max:255',
        ];
    }

    public function messages(): array
    {
        return [
            'name.required' => 'Поле "Назва Дii" є обов\'язковим для заповнення.',
            'name.string' => '"Назва Дii" має бути рядком.',
            'name.max' => '"Назва Дii" не може перевищувати 255 символів.',
            'name.unique' => 'Такий доступ з цим маршрутом вже існує.',
            'route.required' => 'Поле "Маршрут" є обов\'язковим для заповнення.',
            'route.string' => 'Маршрут має бути рядком.',
            'route.max' => 'Маршрут не може перевищувати 255 символів.',
            'route.unique' => 'Такий Маршрут з цим доступом вже існує.',
            'description.string' => 'Опис має бути рядком.',
            'description.max' => 'Опис не може перевищувати 255 символів.',
        ];
    }
}
